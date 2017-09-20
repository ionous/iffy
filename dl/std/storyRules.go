package std

import "github.com/ionous/iffy/spec"

func storyRules(c spec.Block) {
	// print the class name if all else fails
	if c.Cmd("run rule", "commence").Begin() {
		if c.Param("decide").Cmds().Begin() {
			c.Cmd("set text", "story", "status left", "{determine playerSurroundings}")
			c.Cmd("set text", "story", "status right", "{score}/{turnCount}")
			c.Cmd("set text", "story", "status right", "{score}/{turnCount}")
			c.End()
		}
		c.End()
	}
	if c.Cmd("text rule", "player surroundings").Begin() {
		// player location or darkness
		if c.Param("decide").Cmd("buffer").Begin() {
			if c.Cmds().Begin() {
				c.Cmd("determine", c.Cmd("print name", c.Cmd("location of", c.Cmd("player"))))
				c.End()
			}
			c.End()
		}
		c.End()
	}
	if c.Cmd("run rule", "construct status line").Begin() {
		c.Cmd("say", "{story.statusLeft}")
		c.Cmd("say", "{story.statusRight}")
		c.End()
	}

	if c.Cmd("run rule", "print banner text").Begin() {
		// FIX: REALLY SHOULD HAVE ELSE IN TEMPLATE
		if c.Cmd("say").Begin() {
			c.Param("text").Val("{story.name}")
			c.Param("else").Val("Welcome")
			c.End()
		}
		c.End()
	}
	// *PrintBannerText
	// *ConstructStatusLine
}

// func updateScore(c spec.Block) {
// 	// print the class name if all else fails
// 	if c.Cmd("run rule", "update score").Begin() {
// 		if c.Param("decide").Cmds().Begin() {
// 			if c.Cmd("choose").Begin() {
// 				if c.Cmd("if", c.Cmd("get", c.Cmd("get", "@", "story"), "scored")).Begin() {

// 					c.End()
// 				}
// 				c.End()
// 			}
// 		}
// 		// 	c.Cmd("say", c.Cmd("class name", c.Cmd("get", "@", "target")))
// 		// 	c.End()
// 		// }
// 		c.End()
// 	}

// 	// FIX: duplication with end turn
// 	if story.Is("scored") {
// 		score := story.Num("score")
// 		status := fmt.Sprintf("%d/%d", int(score), int(0))
// 		g.The("status bar").SetText("right", status)
// 	}
// 	room := g.The("player").Object("whereabouts")
// 	if !room.Exists() {
// 		rooms := g.List("rooms")
// 		if rooms.Len() == 0 {
// 			panic("story has no rooms")
// 		}
// 		room = rooms.Get(0).Object()
// 	}
// 	story.Go("set initial position", g.The("player"), room)
// 	story.Go("print the banner") // see: banner.go
// 	room = g.The("player").Object("whereabouts")
// 	// FIX: Go() should handle both Name() and ref
// 	story.Go("describe the first room", room)
// 	story.IsNow("playing")
// }))
