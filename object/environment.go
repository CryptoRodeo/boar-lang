package object

type Environment struct {
	store map[string]Object
	outer *Environment //outer scope
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	/**
		If we cant find the identifier in the current scope
		and we have an enclosing, outer scope, search in that scope
	**/
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

/**
dev notes:
- we need to preserve the bindings (let x = 1, let i = fn(){}) while at the same time
making new ones available. i.e. we need to extend the scope / environment

- Extending the scope/env simply means to create a new instance of object.Environment with a pointer to the
one its extending (the outer scope). Doing so allows us to create a fresh, empty scope "within" the previous one

Searching for associations:
- when the new environment's/scope's Get method is called and it doesn't have the value
associated with the given name, it should call the Get method of the enclosing environment (outer scope)
- it should do this until there is no enclosing environment anymore (i.e. we've reached the global scope) and throw an error



**/
