package node

import (
	"context"
	"crypto/rand"
	mrand "math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	abcicli "github.com/tendermint/tendermint/abci/client"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"

	"github.com/celestiaorg/rollmint/config"
	"github.com/celestiaorg/rollmint/da"
	mockda "github.com/celestiaorg/rollmint/da/mock"
	"github.com/celestiaorg/rollmint/mocks"
	"github.com/celestiaorg/rollmint/p2p"
	"github.com/celestiaorg/rollmint/store"
	rmtypes "github.com/celestiaorg/rollmint/types"
)

func TestAggregatorMode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	app := &mocks.Application{}
	app.On("InitChain", mock.Anything).Return(abci.ResponseInitChain{})
	app.On("CheckTx", mock.Anything).Return(abci.ResponseCheckTx{})
	app.On("BeginBlock", mock.Anything).Return(abci.ResponseBeginBlock{})
	app.On("DeliverTx", mock.Anything).Return(abci.ResponseDeliverTx{})
	app.On("EndBlock", mock.Anything).Return(abci.ResponseEndBlock{})
	app.On("Commit", mock.Anything).Return(abci.ResponseCommit{})
	app.On("PrepareProposal", mock.Anything).
		Return(func(req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
			return abci.ResponsePrepareProposal{Txs: req.Txs}
		})

	key, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	signingKey, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	anotherKey, _, _ := crypto.GenerateEd25519Key(rand.Reader)

	blockManagerConfig := config.BlockManagerConfig{
		BlockTime:   1 * time.Second,
		NamespaceID: rmtypes.NamespaceID{1, 2, 3, 4, 5, 6, 7, 8},
	}
	node, err := NewNode(context.Background(), config.NodeConfig{DALayer: "mock", Aggregator: true, BlockManagerConfig: blockManagerConfig}, key, signingKey, abcicli.NewLocalClient(nil, app), &types.GenesisDoc{ChainID: "test"}, log.TestingLogger())
	require.NoError(err)
	require.NotNil(node)

	assert.False(node.IsRunning())

	err = node.Start()
	assert.NoError(err)
	defer func() {
		err := node.Stop()
		assert.NoError(err)
	}()
	assert.True(node.IsRunning())

	pid, err := peer.IDFromPrivateKey(anotherKey)
	require.NoError(err)
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				node.incomingTxCh <- &p2p.GossipMessage{Data: []byte(time.Now().String()), From: pid}
				time.Sleep(time.Duration(mrand.Uint32()%20) * time.Millisecond)
			}
		}
	}()
	time.Sleep(3 * time.Second)
	cancel()
}

// TestTxGossipingAndAggregation setups a network of nodes, with single aggregator and multiple producers.
// Nodes should gossip transactions and aggregator node should produce blocks.
func TestTxGossipingAndAggregation(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var wg sync.WaitGroup
	clientNodes := 4
	aggCtx, aggCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	nodes, apps := createNodes(aggCtx, ctx, clientNodes+1, &wg, t)

	wg.Add((clientNodes + 1) * clientNodes)
	for _, n := range nodes {
		require.NoError(n.Start())
	}

	time.Sleep(1 * time.Second)

	for i := 1; i < len(nodes); i++ {
		data := strconv.Itoa(i) + time.Now().String()
		require.NoError(nodes[i].P2P.GossipTx(context.TODO(), []byte(data)))
	}

	timeout := time.NewTimer(time.Second * 30)
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		wg.Wait()
	}()
	select {
	case <-doneChan:
	case <-timeout.C:
		t.Fatal("failing after timeout")
	}

	require.NoError(nodes[0].Stop())
	aggCancel()
	time.Sleep(100 * time.Millisecond)
	for _, n := range nodes[1:] {
		require.NoError(n.Stop())
	}
	cancel()
	time.Sleep(100 * time.Millisecond)
	aggApp := apps[0]
	apps = apps[1:]

	aggApp.AssertNumberOfCalls(t, "DeliverTx", clientNodes)
	aggApp.AssertExpectations(t)

	for i, app := range apps {
		app.AssertNumberOfCalls(t, "DeliverTx", clientNodes)
		app.AssertExpectations(t)

		// assert that we have most of the blocks from aggregator
		beginCnt := 0
		endCnt := 0
		commitCnt := 0
		for _, call := range app.Calls {
			switch call.Method {
			case "BeginBlock":
				beginCnt++
			case "EndBlock":
				endCnt++
			case "Commit":
				commitCnt++
			}
		}
		aggregatorHeight := nodes[0].Store.Height()
		adjustedHeight := int(aggregatorHeight - 3) // 3 is completely arbitrary
		assert.GreaterOrEqual(beginCnt, adjustedHeight)
		assert.GreaterOrEqual(endCnt, adjustedHeight)
		assert.GreaterOrEqual(commitCnt, adjustedHeight)

		// assert that all blocks known to node are same as produced by aggregator
		for h := uint64(1); h <= nodes[i].Store.Height(); h++ {
			aggBlock, err := nodes[0].Store.LoadBlock(h)
			require.NoError(err)
			nodeBlock, err := nodes[i].Store.LoadBlock(h)
			require.NoError(err)
			assert.Equal(aggBlock, nodeBlock)
		}
	}
}

func createNodes(aggCtx, ctx context.Context, num int, wg *sync.WaitGroup, t *testing.T) ([]*Node, []*mocks.Application) {
	t.Helper()

	if aggCtx == nil {
		aggCtx = context.Background()
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// create keys first, as they are required for P2P connections
	keys := make([]crypto.PrivKey, num)
	for i := 0; i < num; i++ {
		keys[i], _, _ = crypto.GenerateEd25519Key(rand.Reader)
	}

	nodes := make([]*Node, num)
	apps := make([]*mocks.Application, num)
	dalc := &mockda.DataAvailabilityLayerClient{}
	_ = dalc.Init([8]byte{}, nil, store.NewDefaultInMemoryKVStore(), log.TestingLogger())
	_ = dalc.Start()

	nodes[0], apps[0] = createNode(aggCtx, 0, true, dalc, keys, wg, t)
	for i := 1; i < num; i++ {
		nodes[i], apps[i] = createNode(ctx, i, false, dalc, keys, wg, t)
	}

	return nodes, apps
}

func createNode(ctx context.Context, n int, aggregator bool, dalc da.DataAvailabilityLayerClient, keys []crypto.PrivKey, wg *sync.WaitGroup, t *testing.T) (*Node, *mocks.Application) {
	t.Helper()
	require := require.New(t)
	// nodes will listen on consecutive ports on local interface
	// random connections to other nodes will be added
	startPort := 10000
	p2pConfig := config.P2PConfig{
		ListenAddress: "/ip4/127.0.0.1/tcp/" + strconv.Itoa(startPort+n),
	}
	bmConfig := config.BlockManagerConfig{
		BlockTime:   300 * time.Millisecond,
		NamespaceID: rmtypes.NamespaceID{8, 7, 6, 5, 4, 3, 2, 1},
	}
	for i := 0; i < len(keys); i++ {
		if i == n {
			continue
		}
		r := i
		id, err := peer.IDFromPrivateKey(keys[r])
		require.NoError(err)
		p2pConfig.Seeds += "/ip4/127.0.0.1/tcp/" + strconv.Itoa(startPort+r) + "/p2p/" + id.Pretty() + ","
	}
	p2pConfig.Seeds = strings.TrimSuffix(p2pConfig.Seeds, ",")

	app := &mocks.Application{}
	app.On("InitChain", mock.Anything).Return(abci.ResponseInitChain{})
	app.On("CheckTx", mock.Anything).Return(abci.ResponseCheckTx{})
	app.On("BeginBlock", mock.Anything).Return(abci.ResponseBeginBlock{})
	app.On("EndBlock", mock.Anything).Return(abci.ResponseEndBlock{})
	app.On("Commit", mock.Anything).Return(abci.ResponseCommit{})
	app.On("DeliverTx", mock.Anything).Return(abci.ResponseDeliverTx{}).Run(func(args mock.Arguments) {
		wg.Done()
	})

	// note that if we are the aggregator, then we prepare the proposal, 
	// whereas if we are the validator we process the proposal to validate it
	if aggregator {
		app.On("PrepareProposal", mock.Anything).
			Return(func(req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
				return abci.ResponsePrepareProposal{Txs: req.Txs}
			})
	} else {
		app.On("ProcessProposal", mock.Anything).Return(abci.ResponseProcessProposal{Result: abci.ResponseProcessProposal_ACCEPT})
		
		// TODO: should only non aggregators be calling this?
		app.On("ExtendVote", mock.Anything).Return(abci.ResponseExtendVote{})
		app.On("VerifyVoteExtension", mock.Anything).Return(abci.ResponseVerifyVoteExtension{})
	}


	if ctx == nil {
		ctx = context.Background()
	}

	signingKey, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	node, err := NewNode(
		ctx,
		config.NodeConfig{
			P2P:                p2pConfig,
			DALayer:            "mock",
			Aggregator:         aggregator,
			BlockManagerConfig: bmConfig,
		},
		keys[n],
		signingKey,
		abcicli.NewLocalClient(nil, app),
		&types.GenesisDoc{ChainID: "test"},
		log.TestingLogger().With("node", n))
	require.NoError(err)
	require.NotNil(node)

	// use same, common DALC, so nodes can share data
	node.dalc = dalc
	node.blockManager.SetDALC(dalc)

	return node, app
}
