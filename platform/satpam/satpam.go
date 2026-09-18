// Package satpam provides authorization for the application
package satpam

type Authorizable interface {
	Authorize(sub *Attribute, res map[string]*Attribute)
	Name() string
}

type Security struct {
	permissions []Authorizable
	resources   map[string]any
}

func New(perms ...Authorizable) *Security {
	return &Security{
		permissions: perms,
		resources:   make(map[string]any),
	}
}

func (s *Security) AddResource(name string, resource any) *Security {
	s.resources[name] = resource
	return s
}

func (s *Security) buildCarriers(subject any) (*Attribute, map[string]*Attribute) {
	subC := NewAttributeCarrier(subject)
	resMap := make(map[string]*Attribute, len(s.resources))
	for k, v := range s.resources {
		resMap[k] = NewAttributeCarrier(v)
	}
	return subC, resMap
}

func (s *Security) GetAllPermissions(subject any) map[string]bool {
	results := make(map[string]bool, len(s.permissions))
	subC, resMap := s.buildCarriers(subject)

	for _, perm := range s.permissions {
		func(p Authorizable) {
			defer func() {
				if r := recover(); r != nil {
					results[p.Name()] = false
				}
			}()
			p.Authorize(subC, resMap)
			results[p.Name()] = true
		}(perm)
	}

	return results
}

func (s *Security) AuthorizeAllPermissions(subject any) {
	subC, resMap := s.buildCarriers(subject)
	for _, perm := range s.permissions {
		perm.Authorize(subC, resMap)
	}
}

func (s *Security) Can(subject any) bool {
	subC, resMap := s.buildCarriers(subject)
	for _, perm := range s.permissions {
		ok := true
		func(p Authorizable) {
			defer func() {
				if r := recover(); r != nil {
					ok = false
				}
			}()
			p.Authorize(subC, resMap)
		}(perm)
		if !ok {
			return false
		}
	}
	return true
}
