package router

import (
	"errors"
	"strings"
	"sync"

	buffpool "github.com/mitghi/lfpool"
	"github.com/mitghi/protox/containers"
	"github.com/mitghi/protox/messages"
)

type Subs struct {
	sync.RWMutex
	*containers.Radix
	*buffpool.BuffPool
	cache subcache
}

type subscription struct {
	topic   string
	uid     string
	eticket byte
	isLeaf  bool
}

type subinfo struct {
	subs map[string]*subscription
}

// - MARK: Initializers.

func NewSubs(ralloc int) *Subs {
	result := &Subs{
		cache:    subcache{make(map[string]*subcacheline), 0, 0, 0, 0},
		Radix:    containers.NewRadix(ralloc),
		BuffPool: buffpool.NewBuffPool(),
	}
	return result
}

func NewSubsWithBuffer(ralloc int, buff *buffpool.BuffPool) *Subs {
	result := &Subs{
		cache:    subcache{make(map[string]*subcacheline), 0, 0, 0, 0},
		Radix:    containers.NewRadix(ralloc),
		BuffPool: buff,
	}
	return result
}

// - MARK: subinfo section.

func (s *subinfo) insert(uid string, topic string, eticket byte, isLeaf bool) (ret *subscription) {
	sb, ok := s.subs[uid]
	if !ok {
		ret = &subscription{topic, uid, eticket, isLeaf}
		s.subs[uid] = ret
		return ret
	}
	sb.uid = uid
	sb.eticket = eticket
	sb.isLeaf = isLeaf
	return sb
}

func (s *subinfo) remove(uid string) bool {
	sb, ok := s.subs[uid]
	if !ok {
		return false
	}
	if sb == nil {
		return false
	}
	delete(s.subs, uid)
	return true
}

func (s *subinfo) get(uid string) (*subscription, bool) {
	sb, ok := s.subs[uid]
	if !ok {
		return nil, false
	}
	return sb, true
}

// - MARK: Subs section.

func (s *Subs) searchPath(path []byte, sep string, callback func(**containers.RDXNode, [][]byte, int, int)) error {
	var (
		buff  *buffpool.Buffer = s.GetBuffer(len(path))
		paths [][]byte
		plen  int
		curr  **containers.RDXNode
	)
	defer s.ReleaseBuffer(buff)
	paths, err := messages.TopicComponents(path)
	if err != nil {
		return errors.New("router: invalid topic.")
	}
	plen = len(paths) - 1
ML:
	for i, v := range paths {
		if i == 0 {
			var bstr string
			buff.Write(v)
			bstr = buff.String()
			_, tail, word := s.Find(bstr)
			if tail != nil && word == bstr {
				curr = &tail
			} else {
				return SNoSubs
			}
			buff.WriteByte('/')
			bstr = buff.String()
			_, ntail, nword := s.FindFrom(*curr, bstr)
			if ntail == nil {
				return SNoSubs
			}
			if n := (*ntail.Link); n != nil && nword == bstr {
				curr = &n
			} else {
				return SNoSubs
			}
			buff.Reset()
			continue
		}
		if *curr == nil {
			break
		}
		if string(v) != IdentWLCD {
			if vn := s.wldCheck(curr); vn != nil && callback != nil {
				callback(&vn, paths, i, plen)
			}
		}
		buff.WriteByte('/')
		buff.Write(v)
		bstr := buff.String()
		_, tail, word := s.FindFrom(*curr, bstr)
		if tail == nil {
			return SNoSubs
		}
		n := (*tail.Link)
		if n == nil {
			return SNoSubs
		}
		if word != bstr {
			return SNoSubs
		}
		curr = &n
		buff.Reset()
		if i == plen {
			if callback != nil {
				callback(curr, paths, i, plen)
				break ML
			}
		} else {
			buff.WriteString((*curr).Key)
			buff.WriteByte('/')
			_, tail, word = s.FindFrom(*curr, buff.String())
			if tail == nil {
				return SNoSubs
			}
			n = (*tail.Link)
			if n == nil {
				return SNoSubs
			}
			curr = &n
		}
		buff.Reset()
	}

	return nil
}

func (s *Subs) wldCheck(node **containers.RDXNode) *containers.RDXNode {
	_, tail, word := s.FindFrom((*node), TOPICWLCD)
	if tail != nil {
		if n := (*tail.Link); n != nil && word == TOPICWLCD {
			return n
		}
	}
	return nil
}

func (s *Subs) insertPath(path []byte, value interface{}, sep string, callback func(**containers.RDXNode, int, int)) **containers.RDXNode {
	var (
		buff  *buffpool.Buffer = s.GetBuffer(len(path))
		paths [][]byte
		plen  int
		curr  **containers.RDXNode
	)
	defer s.ReleaseBuffer(buff)
	// TODO
	//  check the error code
	paths, _ = messages.TopicComponents(path)
	plen = len(paths) - 1
	for i, v := range paths {
		if i == 0 {
			buff.Write(v)
			a, _ := s.InsertRFrom(s.GetRoot(), buff.String(), value)
			curr = &a
			a.NoMrg = true
			buff.Reset()
			buff.WriteString(a.Key)
			buff.WriteByte('/')
			b, _ := s.InsertRFrom(*curr, buff.String(), value)
			curr = &b
			b.NoMrg = true
			buff.Reset()
			continue
		}
		buff.WriteByte('/')
		buff.Write(v)
		a, _ := s.InsertRFrom(*curr, buff.String(), value)
		if a != nil {
			curr = &a
			a.NoMrg = true
			buff.Reset()
			buff.WriteString(a.Key)
		} else {
			logger.Fatal("(InsertPath) FATAl. TODO: FATAL SITUATION, CHECK.", "a", a)
			/* TODO */
			// NOTE: this event is fatal and should not occur
			// uncomment to crash the program in case of an error
			// and get stack trace dump.
			// panic("Subs: inconsistent state")
		}
		if buff.Len() > 0 {
			lch := buff.Bytes()[buff.Len()-1]
			/* d e b u g */
			// logger.Debug("len(buff)>0", lch, buff.String())
			// logger.Debug("insertPath()", "lch", lch)
			/* d e b u g */
			// TODO
			// . remove this later
			if strings.ContainsAny(string(lch), sep) == true {
			}
			n := *curr
			/* d e b u g */
			// NOTE: CHANGED:
			// if n != nil && n.Key == IdentWLCD {
			/* d e b u g */
			if n != nil {
				if n.Key == IdentWLCD {
					if callback != nil {
						callback(curr, i, plen)
					}
				} else if i == plen {
					/* d e b u g */
					// logger.FInfof("insertPath", "** [Router] in i==plen condition, plen(%d).", plen)
					/* d e b u g */
					if callback != nil {
						callback(curr, i, plen)
					}
				}
				// TODO: NOTE:
				//  this is for using last 1 byte in struct
				//  as option flags.
				/* d e b u g */
				// else {
				// 	n.SetOpts(0x1)
				// }
				/* d e b u g */
			}
		}
		if i != plen {
			buff.WriteByte('/')
			b, _ := s.InsertRFrom(*curr, buff.String(), value)
			curr = &b
			b.NoMrg = true
		}
		buff.Reset()
	}
	return curr
}

func newSubinfo() *subinfo {
	return &subinfo{make(map[string]*subscription)}
}
