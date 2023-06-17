package accounts

type State struct {
	stateTrie *Trie
}

var JJETH_STATE *State

func init() {
	JJETH_STATE = NewState()
}
func NewState() *State {
	state := new(State)
	state.stateTrie = NewTrie()
	return state
}

func (state *State) putAccount(address string, accountData []byte) {
	state.stateTrie.Put(address, accountData)
}

func (state *State) GetAccount(address string) []byte {
	return state.stateTrie.Get(address)
}

func (state *State) GetStateRoot() []byte {
	return state.stateTrie.RootHash
}
