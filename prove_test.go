package golog

import "testing"

func TestFacts (t *testing.T) {
    m := NewMachine().Consult(`
        father(michael).
        father(marc).

        mother(gail).

        parent(X) :-
            father(X).
        parent(X) :-
            mother(X).
    `)

    // these should be provably true
    if !m.CanProve(`father(michael).`) {
        t.Errorf("Couldn't prove father(michael)")
    }
    if !m.CanProve(`father(marc).`) {
        t.Errorf("Couldn't prove father(marc)")
    }
    if !m.CanProve(`parent(michael).`) {
        t.Errorf("Couldn't prove parent(michael)")
    }
    if !m.CanProve(`parent(marc).`) {
        t.Errorf("Couldn't prove parent(marc)")
    }

    // these should not be provable
    if m.CanProve(`father(sue).`) {
        t.Errorf("Proved father(sue)")
    }
    if m.CanProve(`mother(michael).`) {
        t.Errorf("Proved mother(michael)")
    }
    if m.CanProve(`parent(sue).`) {
        t.Errorf("Proved parent(sue)")
    }

    // trivial predicate with multiple solutions
    solutions := m.ProveAll(`father(X).`)
    if len(solutions) != 2 {
        t.Errorf("Wrong number of solutions: %d vs 2", len(solutions))
    }
    if x := solutions[0].ByName_("X").String(); x != "michael" {
        t.Errorf("Wrong first solution: %s", x)
    }
    if x := solutions[1].ByName_("X").String(); x != "marc" {
        t.Errorf("Wrong second solution: %s", x)
    }

    // simple predicate with multiple solutions
    solutions = m.ProveAll(`parent(Name).`)
    if len(solutions) != 3 {
        t.Errorf("Wrong number of solutions: %d vs 2", len(solutions))
    }
    if x := solutions[0].ByName_("Name").String(); x != "michael" {
        t.Errorf("Wrong first solution: %s", x)
    }
    if x := solutions[1].ByName_("Name").String(); x != "marc" {
        t.Errorf("Wrong second solution: %s", x)
    }
    if x := solutions[2].ByName_("Name").String(); x != "gail" {
        t.Errorf("Wrong third solution: %s", x)
    }
}

func TestConjunction(t *testing.T) {
    m := NewMachine().Consult(`
        floor_wax(briwax).
        floor_wax(shimmer).
        floor_wax(minwax).

        dessert(shimmer).
        dessert(cake).
        dessert(pie).

        verb(glimmer).
        verb(shimmer).

        snl(Item) :-
            floor_wax(Item),
            dessert(Item).

        three(Item) :-
            verb(Item),
            dessert(Item),
            floor_wax(Item).
    `)

    skits := m.ProveAll(`snl(X).`)
    if len(skits) != 1 {
        t.Errorf("Wrong number of solutions: %d vs 1", len(skits))
    }
    if x := skits[0].ByName_("X").String(); x != "shimmer" {
        t.Errorf("Wrong solution: %s vs shimmer", x)
    }

    skits = m.ProveAll(`three(W).`)
    if len(skits) != 1 {
        t.Errorf("Wrong number of solutions: %d vs 1", len(skits))
    }
    if x := skits[0].ByName_("W").String(); x != "shimmer" {
        t.Errorf("Wrong solution: %s vs shimmer", x)
    }
}

func TestCut(t *testing.T) {
    m := NewMachine().Consult(`
        single(foo) :-
            !.
        single(bar).

        twice(X) :-
            single(X).  % cut inside here doesn't cut twice/1
        twice(bar).
    `)

    proofs := m.ProveAll(`single(X).`)
    if len(proofs) != 1 {
        t.Errorf("Wrong number of solutions: %d vs 1", len(proofs))
    }
    if x := proofs[0].ByName_("X").String(); x != "foo" {
        t.Errorf("Wrong solution: %s vs foo", x)
    }

    proofs = m.ProveAll(`twice(X).`)
    if len(proofs) != 2 {
        t.Errorf("Wrong number of solutions: %d vs 2", len(proofs))
    }
    if x := proofs[0].ByName_("X").String(); x != "foo" {
        t.Errorf("Wrong solution: %s vs foo", x)
    }
    if x := proofs[1].ByName_("X").String(); x != "bar" {
        t.Errorf("Wrong solution: %s vs bar", x)
    }
}

func TestAppend(t *testing.T) {
    m := NewMachine().Consult(`
        append([], A, A).   % test same variable name as other clauses
        append([A|B], C, [A|D]) :-
            append(B, C, D).
    `)

    proofs := m.ProveAll(`append([a], [b], List).`)
    if len(proofs) != 1 {
        t.Errorf("Wrong number of answers: %d vs 1", len(proofs))
    }
    if x := proofs[0].ByName_("List").String(); x != "'.'(a, '.'(b, []))" {
        t.Errorf("Wrong solution: %s vs '.'(a, '.'(b, []))", x)
    }

    proofs = m.ProveAll(`append([a,b,c], [d,e], List).`)
    if len(proofs) != 1 {
        t.Errorf("Wrong number of answers: %d vs 1", len(proofs))
    }
    if x := proofs[0].ByName_("List").String(); x != "'.'(a, '.'(b, '.'(c, '.'(d, '.'(e, [])))))" {
        t.Errorf("Wrong solution: %s", x)
    }
}

func TestCall (t *testing.T) {
    m := NewMachine().Consult(`
        bug(spider).
        bug(fly).

        squash(Animal, Class) :-
            call(Class, Animal).
    `)

    proofs := m.ProveAll(`squash(It, bug).`)
    if len(proofs) != 2 {
        t.Errorf("Wrong number of answers: %d vs 2", len(proofs))
    }
    if x := proofs[0].ByName_("It").String(); x != "spider" {
        t.Errorf("Wrong solution: %s vs spider", x)
    }
    if x := proofs[1].ByName_("It").String(); x != "fly" {
        t.Errorf("Wrong solution: %s vs fly", x)
    }
}

func TestDisjunction (t *testing.T) {
    m := NewMachine().Consult(`
        insect(fly).
        arachnid(spider).
        squash(Critter) :-
            arachnid(Critter) ; insect(Critter).
    `)

    proofs := m.ProveAll(`squash(It).`)
    if len(proofs) != 2 {
        t.Errorf("Wrong number of answers: %d vs 2", len(proofs))
    }
    if x := proofs[0].ByName_("It").String(); x != "spider" {
        t.Errorf("Wrong solution: %s vs spider", x)
    }
    if x := proofs[1].ByName_("It").String(); x != "fly" {
        t.Errorf("Wrong solution: %s vs fly", x)
    }
}

func BenchmarkTrue(b *testing.B) {
    m := NewMachine()
    for i := 0; i < b.N; i++ {
        _ = m.ProveAll(`true.`)
    }
}

func BenchmarkAppend(b *testing.B) {
    m := NewMachine().Consult(`
        append([], A, A).   % test same variable name as other clauses
        append([A|B], C, [A|D]) :-
            append(B, C, D).
    `)

    for i := 0; i < b.N; i++ {
        _ = m.ProveAll(`append([a,b,c], [d,e], List).`)
    }
}
