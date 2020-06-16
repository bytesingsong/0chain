package conductrpc

import (
	"net/rpc"

	"0chain.net/conductor/config"
)

// Client of the conductor RPC server.
type Client struct {
	address string
	client  *rpc.Client
}

// NewClient creates new client will be interacting
// with server with given address.
func NewClient(address string) (c *Client, err error) {
	if address, err = Host(address); err != nil {
		return
	}
	c = new(Client)
	if c.client, err = rpc.Dial("tcp", address); err != nil {
		return nil, err
	}
	c.address = address
	return
}

func (c *Client) dial() (err error) {
	c.client, err = rpc.Dial("tcp", c.address)
	return
}

// Address of RPC server.
func (c *Client) Address() string {
	return c.address
}

//
// miner SC RPC
//

func (c *Client) Phase(phase *PhaseEvent) (err error) {
	err = c.client.Call("Server.Phase", phase, &struct{}{})
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.Phase", phase, &struct{}{})
	}
	return
}

// ViewChange notification.
func (c *Client) ViewChange(viewChange *ViewChangeEvent) (err error) {
	err = c.client.Call("Server.ViewChange", viewChange, &struct{}{})
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.ViewChange", viewChange, &struct{}{})
	}
	return
}

// AddMiner notification.
func (c *Client) AddMiner(add *AddMinerEvent) (err error) {
	err = c.client.Call("Server.AddMiner", add, &struct{}{})
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.AddMiner", add, &struct{}{})
	}
	return
}

// AddSharder notification.
func (c *Client) AddSharder(add *AddSharderEvent) (err error) {
	err = c.client.Call("Server.AddSharder", add, &struct{}{})
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.AddSharder", add, &struct{}{})
	}
	return
}

// NodeReady notification.
func (c *Client) NodeReady(nodeID NodeID) (join bool, err error) {
	err = c.client.Call("Server.NodeReady", nodeID, &join)
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.NodeReady", nodeID, &join)
	}
	return
}

// Round notification.
func (c *Client) Round(re *RoundEvent) (err error) {
	err = c.client.Call("Server.Round", re, &struct{}{})
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.Round", re, &struct{}{})
	}
	return
}

// ContributeMPK notification.
func (c *Client) ContributeMPK(cmpke *ContributeMPKEvent) (err error) {
	err = c.client.Call("Server.ContributeMPK", cmpke, &struct{}{})
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.ContributeMPK", cmpke, &struct{}{})
	}
	return
}

// ShareOrSignsShares notification.
func (c *Client) ShareOrSignsShares(sosse *ShareOrSignsSharesEvent) (err error) {
	err = c.client.Call("Server.ShareOrSignsShares", sosse, &struct{}{})
	if err == rpc.ErrShutdown {
		if err = c.dial(); err != nil {
			return
		}
		err = c.client.Call("Server.ShareOrSignsShares", sosse, &struct{}{})
	}
	return
}

//
// state
//

// SendShareOnly to configured nodes. The long polling method, e.g. it blocks.
func (c *Client) SendShareOnly(me NodeID) (only []NodeID, err error) {
	err = c.client.Call("Server.SendShareOnly", me, &only)
	for err == rpc.ErrShutdown {
		err = c.client.Call("Server.SendShareOnly", me, &only)
	}
	return
}

// SendShareBad sends bad share to resulting nodes. To send bad shares only to
// X nodes, use SendShareOnly (nil, nil) with SendShareBad (list, nil).
// The long polling method, e.g. it blocks.
func (c *Client) SendShareBad(me NodeID) (bad []NodeID, err error) {
	err = c.client.Call("Server.SendShareBad", me, &bad)
	for err == rpc.ErrShutdown {
		err = c.client.Call("Server.SendShareBad", me, &bad)
	}
	return
}

func (c *Client) IsRevealed(node NodeID) (pin bool, err error) {
	err = c.client.Call("Server.IsRevealed", node, &pin)
	for err == rpc.ErrShutdown {
		err = c.client.Call("Server.IsRevealed", node, &pin)
	}
	return
}

// State is node state.
type State struct {
	IsMonitor  bool // send monitor events (round, phase, etc)
	Lock       bool // node locked
	IsRevealed bool // revealed shares
	// Byzantine state. Below, if a value is nil, then node behaves as usual
	// for it.
	//
	// Byzantine blockchain
	VRFS                        *config.VRFS
	RoundTimeout                *config.RoundTimeout
	CompetingBlock              *config.CompetingBlock
	SignOnlyCompetingBlocks     *config.SignOnlyCompetingBlocks
	DoubleSpendTransaction      *config.DoubleSpendTransaction
	WrongBlockSignHash          *config.WrongBlockSignHash
	WrongBlockSignKey           *config.WrongBlockSignKey
	WrongBlockHash              *config.WrongBlockHash
	VerificationTicket          *config.VerificationTicket
	WrongVerificationTicketHash *config.WrongVerificationTicketHash
	WrongVerificationTicketKey  *config.WrongVerificationTicketKey
	WrongNotarizedBlockHash     *config.WrongNotarizedBlockHash
	WrongNotarizedBlockKey      *config.WrongNotarizedBlockKey
	NotarizeOnlyCompetingBlock  *config.NotarizeOnlyCompetingBlock
	NotarizedBlock              *config.NotarizedBlock
	// Byzantine View Change
	MPK        *config.MPK
	Shares     *config.Shares
	Signatures *config.Signatures
	Publish    *config.Publish
}
