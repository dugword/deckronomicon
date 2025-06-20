# Using Deckronomicon

## Interactive Player Agent

When running in interactive mode, the player agent will prompt you for
actions.

These are the availabe commands

### Standard Actions

#### `activate` `tap`

Activate an ability of a object. You will be prompted to select the object and
the ability or you can specify the object name or ID: `activite <name|id>`

If the ability requries targets you will be prompted with a choice of targets
to select.

#### `cast`

Cast a spell from your hand. You will be prompted to select the spell or you
can specify the spell name or ID: `cast <name|id>`

If the spell requires targets you will be prompted with a choice of targets
to select.

#### `cheat`

Enable cheats, allowing you to use commands which break the rules of the game.
These are useful for test and debugging.

#### `clear`

Clear all revealed cards from the game.

#### `concede` `exit` `quit`

Concede the game, ending it immediately. The game will be recorded as a loss.

#### `help`

Print the in game help text.

#### `pass` `next` `done`

Pass priority to the next player. This progresses the priority loop either
resolving the next spell on the stack or ending the current step.

#### `play`

Play a land from your hand. You will be prompted to select the land or you can
specify the land name or ID: `play <name|id>`

#### `view`

View a game object with additional details. You will be prompted to select the
card or you can specify the card name or ID: `view <name|id>`

### Cheat Commands

Cheat commands are only available when cheats are enabled in the game. They
are primarily used for testing and debugging purposes, but they can also be
useful when exploring new strategies. Enabling and using cheats may break the
replayability of the game log, and prevent analytics tools from working
correctly runs which have them enabled.

#### `addmana`

Add mana to your mana pool. You must specify the amount and color of mana to
add using a mana string: `addmana {U}{U}` to add two blue mana.

#### `conjure`

Conjure a card into your hand. The conjure cheat will load a new card from the
definitions and add it to the current game state. You must specify the card
name.

#### `draw`

Draw the top card from your library.

#### `discard`

Discard a card from your hand. You must specify the card name or ID to
discard: `discard <name|id>`

### `effect`

Apply a effect to the game. You must specify the effect name and any modifiers
to apply:`effect <name> <modifier>`, e.g. `effect Scry {"Count": 2}` to scry
two cards.

#### `find` `tutor`

Find a card in your library and put it into your hand. You will be prompted to
select the card or you can specify the card name or ID: `find <name|id>`

#### `landdrop`

Reset your land drop for the turn, allowing you to play an additional land.

#### `peek`

Peek at the top card of your library without drawing it.

#### `shuffle`

Shuffle your library, randomizing the order of the cards.

#### `untap`

Untap all your tapped permanents. You will be prompted to select the
permanent or you can specify the permanent name or ID: `untap <name|id>`
