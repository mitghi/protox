package server

//-------------------------------------------------------
// D E B U G - S E C T I O N
//-------------------------------------------------------
// TODO: NOTE:
//  the commented code is the debug version and will be removed when `insertPath`
//  passes all test cases.
//
// func (s *Subs) insertPathDEBUG(path []byte, value interface{}, sep string) **containers.RDXNode {
// 	var (
// 		buff  bytes.Buffer
// 		paths [][]byte
// 		plen  int
// 		curr  **containers.RDXNode
// 	)
// 	// TODO
// 	//  check the error code
// 	paths, _ = messages.TopicComponents(path)
// 	plen = len(paths) - 1
// 	fmt.Printf("(insertPath)path info.\n\tpath(%s)\n", string(path))
// 	for i, v := range paths {
// 		if i == 0 {
// 			buff.Write(v)
// 			fmt.Printf("(insertPath) path and buff info (before first curr).\n\tpath(%s)\n\tbuff(%s)\n", string(v), buff.String())
// 			a, ct := s.InsertRFrom(s.GetRoot(), buff.String(), value)
// 			curr = &a
// 			fmt.Println("BEFORE REMOVING BUFFER STUFF:", a.Key, buff.String())
// 			// for i := 0; i < buff.Len()-ct; i++ {
// 			// 	buff.ReadByte()
// 			// }
// 			buff.Reset()
// 			buff.WriteString(a.Key)
// 			// if len(a.Key) < buff.Len() {
// 			// 	for i := 0; i <= len(a.Key); i++ {
// 			// 		buff.ReadByte()
// 			// 	}
// 			// }
// 			fmt.Printf("(insertPath)first curr and buff\n\tcurr(%v)\n\tbuff(%s)(%d)\n", *curr, buff.String(), ct)
// 			buff.WriteByte('/')
// 			b, ct := s.InsertRFrom(*curr, buff.String(), value)
// 			// b := r.InsertRFrom(*curr, "/", value)
// 			curr = &b
// 			fmt.Printf("(insertPath)secondcurr.\n\tcurr(%v)(%d)\nbuff(%s)\n", *curr, ct, buff.String())
// 			buff.Reset()
// 			continue
// 		}
// 		buff.WriteByte('/')
// 		buff.Write(v)
// 		fmt.Printf("(insertPath) buff is.\n\tbuff(%s)\n", buff.String())
// 		a, ct := s.InsertRFrom(*curr, buff.String(), value)
// 		fmt.Printf("(insertPath) AFTER buff isAFTER.\n\tbuff(%s)\na(%v)\nct(%d)", buff.String(), *a, ct)

// 		if a != nil {
// 			curr = &a
// 			buff.Reset()
// 			buff.WriteString(a.Key)
// 		} else {
// 			fmt.Println("INSDIE THE MAIN LOOP AND FIND OUT THAT A IS NIL", a)
// 		}
// 		// fmt.Printf("(insertPath) b info.\n\tb(%s)(%d)\n", string(b), ct)
// 		if buff.Len() > 0 {
// 			lch := buff.Bytes()[buff.Len()-1]
// 			fmt.Printf("(insertPath)curr(main loop) and lch info.\n\tcurr(%v)\n\tlch(%s)\nbuffer(%s)\n", *curr, string(lch), buff.String())
// 			if strings.ContainsAny(string(lch), sep) == true {
// 				n := *curr
// 				if n != nil && n.Key == "*" {
// 					n.SetOpts(0x1)
// 				}
// 			} else {
// 				fmt.Printf("(insertPath) no wildcard.\n")
// 			}
// 		}
// 		if i != plen {
// 			buff.WriteByte('/')
// 			b, ct := s.InsertRFrom(*curr, buff.String(), value)
// 			curr = &b

// 			fmt.Printf("(insertPath) CURRB HAS FOLLOWING LINK.\n\tlink(%v)(%d)\n", *curr, ct)
// 		}
// 		buff.Reset()
// 	}
// 	return curr
// }

// ---------------------------------------
//
//
//
// func (s *Subs) searchPath(path []byte, sep string, callback func(**containers.RDXNode, [][]byte, int)) error {
// 	var (
// 		buff  bytes.Buffer
// 		paths [][]byte
// 		plen  int
// 		curr  **containers.RDXNode
// 	)
// 	paths, err := messages.TopicComponents(path)
// 	if err != nil {
// 		return errors.New("router: invalid topic.")
// 	}
// 	plen = len(paths) - 1
// ML:
// 	for i, v := range paths {
// 		if i == 0 {
// 			buff.Write(v)
// 			head, tail, word := s.Find(buff.String())
// 			if tail != nil && word == buff.String() {
// 				fmt.Println("works", head, tail, word)
// 				curr = &tail
// 			} else {
// 				fmt.Println("searchPath unable to locate.")
// 				return errors.New("searchPath: cannot locate.")
// 			}
// 			buff.WriteByte('/')
// 			nhead, ntail, nword := s.FindFrom(*curr, buff.String())
// 			if ntail == nil {
// 				fmt.Println("searchPath unable to go beyond 0 iteration..0")
// 				return errors.New("searchPath: cannot locate.")
// 			}
// 			if n := (*ntail.Link); n != nil {
// 				fmt.Println("NWORD:", nword, nhead, ntail, n)
// 				curr = &n
// 			} else {
// 				fmt.Println("searchPath unable to go beyond 0 iteration")
// 			}
// 			buff.Reset()
// 			continue
// 		}
// 		if *curr == nil {
// 			break
// 		}

// 		fmt.Println("HAS WILDCARD?: ")
// 		if string(v) != "*" {
// 			if vn := s.wldCheck(curr); vn != nil && callback != nil {
// 				callback(&vn, paths, i)
// 			}
// 		}

// 		buff.WriteByte('/')
// 		buff.Write(v)
// 		head, tail, word := s.FindFrom(*curr, buff.String())
// 		if tail == nil {
// 			fmt.Println("un1, buff; ", buff.String(), word, tail, curr, *curr, (*curr).Link, (*curr).Next)
// 			return errors.New("searchPath: unable to locate1")
// 		}
// 		fmt.Println("BUFF 1:", buff.String())
// 		n := (*tail.Link)
// 		if n == nil {
// 			return errors.New("searchPath: unable to locate2")
// 		}
// 		if word != buff.String() {
// 			fmt.Println("word!=buff", word, buff.String(), *n, head)
// 			return errors.New("searchPath: unable to locate3")
// 		}
// 		curr = &n
// 		fmt.Println(i, v, "iter", *n)
// 		buff.Reset()

// 		if i == plen {
// 			if callback != nil {
// 				fmt.Println("calling callback ( i==plen )")
// 				callback(curr, paths, i)
// 				break ML
// 			}
// 		} else {
// 			buff.WriteString((*curr).Key)
// 			buff.WriteByte('/')
// 			fmt.Println("BUFF 2:", buff.String())
// 			head, tail, word = s.FindFrom(*curr, buff.String())
// 			if tail == nil {
// 				return errors.New("searchPath: unable to locate4")
// 			}
// 			n = (*tail.Link)
// 			if n == nil {
// 				return errors.New("searchPath: unable to locate5")
// 			}
// 			curr = &n
// 			fmt.Println("not last iter", word, *curr)
// 		}
// 		buff.Reset()
// 	}
// 	fmt.Println("(searchPath) post break.")
// 	return nil
// }

/*                                D E B U G                                        */

// func (s *Server) NotifySubscribe(topic string, prc protobase.ProtoConnection) {
// 	log.Println("+ [Client][Layer] Attached to stream.")
// 	s.Lock()
// 	defer s.Unlock()
// 	_, ok := s.Router[topic]
// 	if !ok {
// 		var cl protobase.ClientInterface = prc.GetClient()
// 		s.Router[topic] = make(map[string]protobase.ProtoConnection)
// 		s.Router[topic][cl.GetIdentifier()] = prc
// 		return
// 	}
// 	var cl protobase.ClientInterface = prc.GetClient()
// 	_, ok = s.Router[topic][cl.GetIdentifier()]
// 	if !ok {
// 		s.Router[topic][cl.GetIdentifier()] = prc
// 		return
// 	}
// }

// func (s *Server) NotifySubscribe(topic string, prc protobase.ProtoConnection)

// TODO add new implementation
// func (s *Server) NotifySubscribe(topic string, prc protobase.ProtoConnection) {
// 	log.Println("+ [Client][Layer] Attached to stream.")
// }

//
// func (s *Server) NotifyPublish(topic string, message string, prc protobase.ProtoConnection) {
// 	s.RLock()
// 	defer s.RUnlock()
// 	_, ok := s.Router[topic]
// 	if !ok {
// 		log.Println("- [Inconsistent] Will.")
// 		return
// 	}
// 	log.Println("** [Service] Directing streams to [Client(s)]")
// 	for _, v := range s.Router[topic] {
// 		v.SendMessage(message, topic)
// 	}

// 	return
// }

// func (s *Server) NotifyConnected(prc protobase.ProtoConnection) {
// 	const fn = "NotifyDisconnected"
// 	s.Lock()
// 	logger.FDebug(fn, "Passed (Genesis) state and is now (Online).", "userid", prc.GetClient().GetIdentifier(), "ip", prc.GetConnection().RemoteAddr().String())
// 	s.Clients[prc.GetConnection()] = prc
// 	s.Unlock()
// 	cl := prc.GetClient()
// 	if c := s.State.get(cl.GetIdentifier()); c != nil {
// 		// fmt.Println("CLIENT ALREADY EXISTS")
// 		logger.FDebug(fn, "client already exists", "client", c)
// 		c.client = cl
// 		c.proto = prc
// 		conn := prc.GetConnection()
// 		c.conn = &conn
// 		c.start = time.Now()
// 	} else {
// 		conn := prc.GetConnection()
// 		c := &connection{time.Now(), &conn, prc, cl}
// 		s.State.set(cl.GetIdentifier(), c)
// 	}
// }

// func (s *Server) NotifyPublish(topic string, message string, prc protobase.ProtoConnection, dir protobase.MsgDir) {
// 	const fn = "NotifyPublish"
// 	m, _ := s.rt.FindSubscribers(topic)
// 	for k, _ := range m {
// 		cl := s.State.get(k)
// 		logger.FDebug(fn, "got the client", cl)
// 		if cl != nil {
// 			var (
// 				user protobase.ClientInterface = cl.proto.GetClient()
// 				clid string                    = cl.uid
// 			)
// 			if cl.proto == prc {
// 				logger.FWarn(fn, "cl.proto == prc", "userId", clid)
// 			}
// 			if stat := cl.proto.GetStatus(); stat == protobase.STATONLINE {
// 				logger.FDebug(fn, "client is online and sending message .... .", "userId", clid, "topic", topic)
// 				// cl.proto.SendMessage(message, topic)

// 				user.Publish(topic, message, protobase.MDOutbound)
// 			} else {
// 				// TODO
// 				// . add to outbound messages
// 			}
// 		}
// 	}
// 	return
// }

// func (s *Server) NotifySubscribe(topic string, prc protobase.ProtoConnection) {
// 	// TODO
// 	const _fn = "NotifySubscribe"
// 	cid := prc.GetClient().GetIdentifier()
// 	logger.FDebug(_fn, "+ [Client][Layer] Attached to stream.", "userId", cid, "stream", topic)
// 	s.rt.AddSub(cid, topic)
// }

// func (r *Router) FindSub(topic string) (map[string]byte, error) {
// 	var (
// 		m map[string]byte = make(map[string]byte)
// 	)

// 	r.RLock()
// 	r.subs.Lock()
// 	defer r.subs.Unlock()
// 	defer r.RUnlock()

// 	callback := func(node **containers.RDXNode, paths [][]byte, level int, plen int) {
// 		n := (*node)
// 		if n.Value != nil {
// 			s := n.Value.(*subinfo).subs
// 			for k, v := range s {
// 				if strs.Match(v.topic, topic, Sep, Wlcd) {
// 					cqos, ok := m[k]
// 					if !ok {
// 						m[k] = v.qos
// 					} else {
// 						if v.qos > cqos {
// 							m[k] = v.qos
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	err := r.subs.searchPath([]byte(topic), Sep, callback)
// 	return m, err
// }

/*                                D E B U G                                        */

// func (s *serverState) setConn(conn net.Conn, info *connection) {
// 	self.Lock()
// 	defer self.Unlock()
// 	self.conns[conn] = info
// 	return nil
// }

// func (s *serverState) getConn(conn net.Conn) (val *connection) {
// 	self.RLock()
// 	defer self.RUnlock()
// 	if val, ok := self.conns[conn]; ok {
// 		return val
// 	}
// 	return nil
// }

// func (s *serverState) getAny(cid interface{}) (val *connection) {
// 	switch cid.(type) {
// 	case string:
// 		s.RLock()
// 		if val, ok := s.clients[cid.(string)]; ok {
// 			s.RUnlock()
// 			return val
// 		}
// 		s.RUnlock()
// 	case net.Conn:
// 		s.RLock()
// 		if val, ok := s.conns[cid.(net.Conn)]; ok {
// 			s.RUnlock()
// 			return val
// 		}
// 		s.RUnlock()
// 	}

// 	return nil
// }

// func (s *serverState) pruneByConn(conn net.Conn) {
// 	self.Lock()
// 	defer self.Unlock()
// 	if val, ok := self.conns[conn]; ok {
// 		delete(self.clients, conn)
// 		return val
// 	}
// 	return nil
// }
