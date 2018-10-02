package auth

/* d e b u g */
// 	// _, isAble := p.abilities[ability]
// 	// _, hasAction := p.actions[action]
// 	// _, hasResource := p.resources[resource]

// 	// return (isAble && hasAction && hasResource)
// }

// func (e *Entity) Add(ability string, action string, resource string) (string, *Action, *Resource) {
// 	// TODO
// 	// . return error in case any of (action, ability or resource) exist.
// 	var (
// 		act *Action
// 		rsc *Resource
// 	)
// 	act, hasAction := e.actions[action]
// 	rsc, hasResource := e.resources[resource]
// 	_, hasAbility := e.abilities[ability]

// 	if !hasAbility {
// 		e.abilities[ability] = ability
// 	}
// 	if !hasAction {
// 		act = &Action{Name: action}
// 		e.actions[action] = act
// 	}
// 	if !hasResource {
// 		rsc = &Resource{Name: resource}
// 		e.resources[resource] = rsc
// 	} else {
// 		rsc.abilities = append(rsc.abilities, ability)
// 		rsc.actions = append(rsc.actions, act)
// 	}
// 	return ability, act, rsc
// }

// func (p *Permission) Add(ability string, action *Action, resource *Resource) {
// 	p.abilities[ability] = ability
// 	p.actions[action.Name] = action
// 	p.resources[resource.Name] = resource
// }

// func (p *Permission) Remove(ability string, action *Action, resource *Resource) bool {
// 	fullrm := true
// 	if _, ok := p.abilities[ability]; ok {
// 		delete(p.abilities, ability)
// 	} else {
// 		fullrm = false
// 	}
// 	if _, ok := p.actions[action.Name]; ok {
// 		delete(p.actions, action.Name)
// 	} else {
// 		fullrm = false
// 	}
// 	if _, ok := p.resources[resource.Name]; ok {
// 		delete(p.resources, resource.Name)
// 	} else {
// 		fullrm = false
// 	}
// 	return fullrm
// }

// func (p *Permission) RemoveByResource(resource *Resource) bool {
// 	if _, ok := p.resources[resource.Name]; ok {
// 		delete(p.resources, resource.Name)
// 		if len(p.resources) == 0 {
// 			p.actions = make(map[string]*Action)
// 			p.abilities = make(map[string]string)
// 		}
// 		return true
// 	}
// 	return false
// }

// func (e *Entity) Add(ability string, action string, resource string) (string, *Action, *Resource) {
// 	// TODO
// 	// . return error in case any of (action, ability or resource) exist.
// 	var (
// 		act *Action
// 		rsc *Resource
// 	)
// 	act, hasAction := e.actions[action]
// 	rsc, hasResource := e.resources[resource]
// 	_, hasAbility := e.abilities[ability]

// 	if !hasAbility {
// 		e.abilities[ability] = ability
// 	}
// 	if !hasAction {
// 		act = &Action{Name: action}
// 		e.actions[action] = act
// 	}
// 	if !hasResource {
// 		rsc = &Resource{Name: resource}
// 		e.resources[resource] = rsc
// 	} else {
// 		rsc.abilities = append(rsc.abilities, ability)
// 		rsc.actions = append(rsc.actions, act)
// 	}
// 	return ability, act, rsc
// }

// func (e *Entity) Add(ability string, action string, resource string) (string, *Action, *Resource) {
// 	// TODO
// 	// . return error in case any of (action, ability or resource) exist.
// 	var (
// 		act *Action
// 		rsc *Resource
// 	)
// 	act, hasAction := e.actions[action]
// 	rsc, hasResource := e.resources[resource]
// 	_, hasAbility := e.abilities[ability]

// 	if !hasAbility {
// 		e.abilities[ability] = ability
// 	}
// 	if !hasAction {
// 		act = &Action{Name: action}
// 		e.actions[action] = act
// 	}
// 	if !hasResource {
// 		rsc = &Resource{Name: resource}
// 		e.resources[resource] = rsc
// 	} else {
// 		rsc.abilities = append(rsc.abilities, ability)
// 		rsc.actions = append(rsc.actions, act)
// 	}
// 	return ability, act, rsc
// }
/* d e b u g */
