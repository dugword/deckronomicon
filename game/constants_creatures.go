package game

// Creature Subtypes
const (
	// Two word subtypes
	SubtypeTimeLord Subtype = "Time Lord"
	// Standard subtypes
	SubtypeAdvisor        Subtype = "Advisor"
	SubtypeAetherborn     Subtype = "Aetherborn"
	SubtypeAlien          Subtype = "Alien"
	SubtypeAlly           Subtype = "Ally"
	SubtypeAngel          Subtype = "Angel"
	SubtypeAntelope       Subtype = "Antelope"
	SubtypeApe            Subtype = "Ape"
	SubtypeArcher         Subtype = "Archer"
	SubtypeArchon         Subtype = "Archon"
	SubtypeArmadillo      Subtype = "Armadillo"
	SubtypeArmy           Subtype = "Army"
	SubtypeArtificer      Subtype = "Artificer"
	SubtypeAssassin       Subtype = "Assassin"
	SubtypeAssemblyWorker Subtype = "Assembly-Worker"
	SubtypeAstartes       Subtype = "Astartes"
	SubtypeAtog           Subtype = "Atog"
	SubtypeAurochs        Subtype = "Aurochs"
	SubtypeAvatar         Subtype = "Avatar"
	SubtypeAzra           Subtype = "Azra"
	SubtypeBadger         Subtype = "Badger"
	SubtypeBalloon        Subtype = "Balloon"
	SubtypeBarbarian      Subtype = "Barbarian"
	SubtypeBard           Subtype = "Bard"
	SubtypeBasilisk       Subtype = "Basilisk"
	SubtypeBat            Subtype = "Bat"
	SubtypeBear           Subtype = "Bear"
	SubtypeBeast          Subtype = "Beast"
	SubtypeBeaver         Subtype = "Beaver"
	SubtypeBeeble         Subtype = "Beeble"
	SubtypeBeholder       Subtype = "Beholder"
	SubtypeBerserker      Subtype = "Berserker"
	SubtypeBird           Subtype = "Bird"
	SubtypeBlinkmoth      Subtype = "Blinkmoth"
	SubtypeBoar           Subtype = "Boar"
	SubtypeBringer        Subtype = "Bringer"
	SubtypeBrushwagg      Subtype = "Brushwagg"
	SubtypeCamarid        Subtype = "Camarid"
	SubtypeCamel          Subtype = "Camel"
	SubtypeCapybara       Subtype = "Capybara"
	SubtypeCaribou        Subtype = "Caribou"
	SubtypeCarrier        Subtype = "Carrier"
	SubtypeCat            Subtype = "Cat"
	SubtypeCentaur        Subtype = "Centaur"
	SubtypeChild          Subtype = "Child"
	SubtypeChimera        Subtype = "Chimera"
	SubtypeCitizen        Subtype = "Citizen"
	SubtypeCleric         Subtype = "Cleric"
	SubtypeClown          Subtype = "Clown"
	SubtypeCockatrice     Subtype = "Cockatrice"
	SubtypeConstruct      Subtype = "Construct"
	SubtypeCoward         Subtype = "Coward"
	SubtypeCoyote         Subtype = "Coyote"
	SubtypeCrab           Subtype = "Crab"
	SubtypeCrocodile      Subtype = "Crocodile"
	SubtypeCTan           Subtype = "C’tan"
	SubtypeCustodes       Subtype = "Custodes"
	SubtypeCyberman       Subtype = "Cyberman"
	SubtypeCyclops        Subtype = "Cyclops"
	SubtypeDalek          Subtype = "Dalek"
	SubtypeDauthi         Subtype = "Dauthi"
	SubtypeDemigod        Subtype = "Demigod"
	SubtypeDemon          Subtype = "Demon"
	SubtypeDeserter       Subtype = "Deserter"
	SubtypeDetective      Subtype = "Detective"
	SubtypeDevil          Subtype = "Devil"
	SubtypeDinosaur       Subtype = "Dinosaur"
	SubtypeDjinn          Subtype = "Djinn"
	SubtypeDoctor         Subtype = "Doctor"
	SubtypeDog            Subtype = "Dog"
	SubtypeDragon         Subtype = "Dragon"
	SubtypeDrake          Subtype = "Drake"
	SubtypeDreadnought    Subtype = "Dreadnought"
	SubtypeDrone          Subtype = "Drone"
	SubtypeDruid          Subtype = "Druid"
	SubtypeDryad          Subtype = "Dryad"
	SubtypeDwarf          Subtype = "Dwarf"
	SubtypeEfreet         Subtype = "Efreet"
	SubtypeEgg            Subtype = "Egg"
	SubtypeElder          Subtype = "Elder"
	SubtypeEldrazi        Subtype = "Eldrazi"
	SubtypeElemental      Subtype = "Elemental"
	SubtypeElephant       Subtype = "Elephant"
	SubtypeElf            Subtype = "Elf"
	SubtypeElk            Subtype = "Elk"
	SubtypeEmployee       Subtype = "Employee"
	SubtypeEye            Subtype = "Eye"
	SubtypeFaerie         Subtype = "Faerie"
	SubtypeFerret         Subtype = "Ferret"
	SubtypeFish           Subtype = "Fish"
	SubtypeFlagbearer     Subtype = "Flagbearer"
	SubtypeFox            Subtype = "Fox"
	SubtypeFractal        Subtype = "Fractal"
	SubtypeFrog           Subtype = "Frog"
	SubtypeFungus         Subtype = "Fungus"
	SubtypeGamer          Subtype = "Gamer"
	SubtypeGargoyle       Subtype = "Gargoyle"
	SubtypeGerm           Subtype = "Germ"
	SubtypeGiant          Subtype = "Giant"
	SubtypeGith           Subtype = "Gith"
	SubtypeGlimmer        Subtype = "Glimmer"
	SubtypeGnoll          Subtype = "Gnoll"
	SubtypeGnome          Subtype = "Gnome"
	SubtypeGoat           Subtype = "Goat"
	SubtypeGoblin         Subtype = "Goblin"
	SubtypeGod            Subtype = "God"
	SubtypeGolem          Subtype = "Golem"
	SubtypeGorgon         Subtype = "Gorgon"
	SubtypeGraveborn      Subtype = "Graveborn"
	SubtypeGremlin        Subtype = "Gremlin"
	SubtypeGriffin        Subtype = "Griffin"
	SubtypeGuest          Subtype = "Guest"
	SubtypeHag            Subtype = "Hag"
	SubtypeHalfling       Subtype = "Halfling"
	SubtypeHamster        Subtype = "Hamster"
	SubtypeHarpy          Subtype = "Harpy"
	SubtypeHellion        Subtype = "Hellion"
	SubtypeHippo          Subtype = "Hippo"
	SubtypeHippogriff     Subtype = "Hippogriff"
	SubtypeHomarid        Subtype = "Homarid"
	SubtypeHomunculus     Subtype = "Homunculus"
	SubtypeHorror         Subtype = "Horror"
	SubtypeHorse          Subtype = "Horse"
	SubtypeHuman          Subtype = "Human"
	SubtypeHydra          Subtype = "Hydra"
	SubtypeHyena          Subtype = "Hyena"
	SubtypeIllusion       Subtype = "Illusion"
	SubtypeImp            Subtype = "Imp"
	SubtypeIncarnation    Subtype = "Incarnation"
	SubtypeInkling        Subtype = "Inkling"
	SubtypeInquisitor     Subtype = "Inquisitor"
	SubtypeInsect         Subtype = "Insect"
	SubtypeJackal         Subtype = "Jackal"
	SubtypeJellyfish      Subtype = "Jellyfish"
	SubtypeJuggernaut     Subtype = "Juggernaut"
	SubtypeKavu           Subtype = "Kavu"
	SubtypeKirin          Subtype = "Kirin"
	SubtypeKithkin        Subtype = "Kithkin"
	SubtypeKnight         Subtype = "Knight"
	SubtypeKobold         Subtype = "Kobold"
	SubtypeKor            Subtype = "Kor"
	SubtypeKraken         Subtype = "Kraken"
	SubtypeLlama          Subtype = "Llama"
	SubtypeLamia          Subtype = "Lamia"
	SubtypeLammasu        Subtype = "Lammasu"
	SubtypeLeech          Subtype = "Leech"
	SubtypeLeviathan      Subtype = "Leviathan"
	SubtypeLhurgoyf       Subtype = "Lhurgoyf"
	SubtypeLicid          Subtype = "Licid"
	SubtypeLizard         Subtype = "Lizard"
	SubtypeManticore      Subtype = "Manticore"
	SubtypeMasticore      Subtype = "Masticore"
	SubtypeMercenary      Subtype = "Mercenary"
	SubtypeMerfolk        Subtype = "Merfolk"
	SubtypeMetathran      Subtype = "Metathran"
	SubtypeMinion         Subtype = "Minion"
	SubtypeMinotaur       Subtype = "Minotaur"
	SubtypeMite           Subtype = "Mite"
	SubtypeMole           Subtype = "Mole"
	SubtypeMonger         Subtype = "Monger"
	SubtypeMongoose       Subtype = "Mongoose"
	SubtypeMonk           Subtype = "Monk"
	SubtypeMonkey         Subtype = "Monkey"
	SubtypeMoonfolk       Subtype = "Moonfolk"
	SubtypeMount          Subtype = "Mount"
	SubtypeMouse          Subtype = "Mouse"
	SubtypeMutant         Subtype = "Mutant"
	SubtypeMyr            Subtype = "Myr"
	SubtypeMystic         Subtype = "Mystic"
	SubtypeNautilus       Subtype = "Nautilus"
	SubtypeNecron         Subtype = "Necron"
	SubtypeNephilim       Subtype = "Nephilim"
	SubtypeNightmare      Subtype = "Nightmare"
	SubtypeNightstalker   Subtype = "Nightstalker"
	SubtypeNinja          Subtype = "Ninja"
	SubtypeNoble          Subtype = "Noble"
	SubtypeNoggle         Subtype = "Noggle"
	SubtypeNomad          Subtype = "Nomad"
	SubtypeNymph          Subtype = "Nymph"
	SubtypeOctopus        Subtype = "Octopus"
	SubtypeOgre           Subtype = "Ogre"
	SubtypeOoze           Subtype = "Ooze"
	SubtypeOrb            Subtype = "Orb"
	SubtypeOrc            Subtype = "Orc"
	SubtypeOrgg           Subtype = "Orgg"
	SubtypeOtter          Subtype = "Otter"
	SubtypeOuphe          Subtype = "Ouphe"
	SubtypeOx             Subtype = "Ox"
	SubtypeOyster         Subtype = "Oyster"
	SubtypePangolin       Subtype = "Pangolin"
	SubtypePeasant        Subtype = "Peasant"
	SubtypePegasus        Subtype = "Pegasus"
	SubtypePentavite      Subtype = "Pentavite"
	SubtypePerformer      Subtype = "Performer"
	SubtypePest           Subtype = "Pest"
	SubtypePhelddagrif    Subtype = "Phelddagrif"
	SubtypePhoenix        Subtype = "Phoenix"
	SubtypePhyrexian      Subtype = "Phyrexian"
	SubtypePilot          Subtype = "Pilot"
	SubtypePincher        Subtype = "Pincher"
	SubtypePirate         Subtype = "Pirate"
	SubtypePlant          Subtype = "Plant"
	SubtypePorcupine      Subtype = "Porcupine"
	SubtypePossum         Subtype = "Possum"
	SubtypePraetor        Subtype = "Praetor"
	SubtypePrimarch       Subtype = "Primarch"
	SubtypePrism          Subtype = "Prism"
	SubtypeProcessor      Subtype = "Processor"
	SubtypeRabbit         Subtype = "Rabbit"
	SubtypeRaccoon        Subtype = "Raccoon"
	SubtypeRanger         Subtype = "Ranger"
	SubtypeRat            Subtype = "Rat"
	SubtypeRebel          Subtype = "Rebel"
	SubtypeReflection     Subtype = "Reflection"
	SubtypeRhino          Subtype = "Rhino"
	SubtypeRigger         Subtype = "Rigger"
	SubtypeRobot          Subtype = "Robot"
	SubtypeRogue          Subtype = "Rogue"
	SubtypeSable          Subtype = "Sable"
	SubtypeSalamander     Subtype = "Salamander"
	SubtypeSamurai        Subtype = "Samurai"
	SubtypeSand           Subtype = "Sand"
	SubtypeSaproling      Subtype = "Saproling"
	SubtypeSatyr          Subtype = "Satyr"
	SubtypeScarecrow      Subtype = "Scarecrow"
	SubtypeScientist      Subtype = "Scientist"
	SubtypeScion          Subtype = "Scion"
	SubtypeScorpion       Subtype = "Scorpion"
	SubtypeScout          Subtype = "Scout"
	SubtypeSculpture      Subtype = "Sculpture"
	SubtypeSerf           Subtype = "Serf"
	SubtypeSerpent        Subtype = "Serpent"
	SubtypeServo          Subtype = "Servo"
	SubtypeShade          Subtype = "Shade"
	SubtypeShaman         Subtype = "Shaman"
	SubtypeShapeshifter   Subtype = "Shapeshifter"
	SubtypeShark          Subtype = "Shark"
	SubtypeSheep          Subtype = "Sheep"
	SubtypeSiren          Subtype = "Siren"
	SubtypeSkeleton       Subtype = "Skeleton"
	SubtypeSkunk          Subtype = "Skunk"
	SubtypeSlith          Subtype = "Slith"
	SubtypeSliver         Subtype = "Sliver"
	SubtypeSloth          Subtype = "Sloth"
	SubtypeSlug           Subtype = "Slug"
	SubtypeSnail          Subtype = "Snail"
	SubtypeSnake          Subtype = "Snake"
	SubtypeSoldier        Subtype = "Soldier"
	SubtypeSoltari        Subtype = "Soltari"
	SubtypeSpawn          Subtype = "Spawn"
	SubtypeSpecter        Subtype = "Specter"
	SubtypeSpellshaper    Subtype = "Spellshaper"
	SubtypeSphinx         Subtype = "Sphinx"
	SubtypeSpider         Subtype = "Spider"
	SubtypeSpike          Subtype = "Spike"
	SubtypeSpirit         Subtype = "Spirit"
	SubtypeSplinter       Subtype = "Splinter"
	SubtypeSponge         Subtype = "Sponge"
	SubtypeSquid          Subtype = "Squid"
	SubtypeSquirrel       Subtype = "Squirrel"
	SubtypeStarfish       Subtype = "Starfish"
	SubtypeSurrakar       Subtype = "Surrakar"
	SubtypeSurvivor       Subtype = "Survivor"
	SubtypeSynth          Subtype = "Synth"
	SubtypeTentacle       Subtype = "Tentacle"
	SubtypeTetravite      Subtype = "Tetravite"
	SubtypeThalakos       Subtype = "Thalakos"
	SubtypeThopter        Subtype = "Thopter"
	SubtypeThrull         Subtype = "Thrull"
	SubtypeTiefling       Subtype = "Tiefling"
	SubtypeToy            Subtype = "Toy"
	SubtypeTreefolk       Subtype = "Treefolk"
	SubtypeTrilobite      Subtype = "Trilobite"
	SubtypeTriskelavite   Subtype = "Triskelavite"
	SubtypeTroll          Subtype = "Troll"
	SubtypeTurtle         Subtype = "Turtle"
	SubtypeTyranid        Subtype = "Tyranid"
	SubtypeUnicorn        Subtype = "Unicorn"
	SubtypeVampire        Subtype = "Vampire"
	SubtypeVarmint        Subtype = "Varmint"
	SubtypeVedalken       Subtype = "Vedalken"
	SubtypeVolver         Subtype = "Volver"
	SubtypeWall           Subtype = "Wall"
	SubtypeWalrus         Subtype = "Walrus"
	SubtypeWarlock        Subtype = "Warlock"
	SubtypeWarrior        Subtype = "Warrior"
	SubtypeWeasel         Subtype = "Weasel"
	SubtypeWeird          Subtype = "Weird"
	SubtypeWerewolf       Subtype = "Werewolf"
	SubtypeWhale          Subtype = "Whale"
	SubtypeWizard         Subtype = "Wizard"
	SubtypeWolf           Subtype = "Wolf"
	SubtypeWolverine      Subtype = "Wolverine"
	SubtypeWombat         Subtype = "Wombat"
	SubtypeWorm           Subtype = "Worm"
	SubtypeWraith         Subtype = "Wraith"
	SubtypeWurm           Subtype = "Wurm"
	SubtypeYeti           Subtype = "Yeti"
	SubtypeZombie         Subtype = "Zombie"
	SubtypeZubera         Subtype = "Zubera"
)

var StringToCreatureSubtype = map[string]Subtype{
	// Two word subtypes
	"Time Lord": SubtypeTimeLord,
	// Standard subtypes
	"Advisor":         SubtypeAdvisor,
	"Aetherborn":      SubtypeAetherborn,
	"Alien":           SubtypeAlien,
	"Ally":            SubtypeAlly,
	"Angel":           SubtypeAngel,
	"Antelope":        SubtypeAntelope,
	"Ape":             SubtypeApe,
	"Archer":          SubtypeArcher,
	"Archon":          SubtypeArchon,
	"Armadillo":       SubtypeArmadillo,
	"Army":            SubtypeArmy,
	"Artificer":       SubtypeArtificer,
	"Assassin":        SubtypeAssassin,
	"Assembly-Worker": SubtypeAssemblyWorker,
	"Astartes":        SubtypeAstartes,
	"Atog":            SubtypeAtog,
	"Aurochs":         SubtypeAurochs,
	"Avatar":          SubtypeAvatar,
	"Azra":            SubtypeAzra,
	"Badger":          SubtypeBadger,
	"Balloon":         SubtypeBalloon,
	"Barbarian":       SubtypeBarbarian,
	"Bard":            SubtypeBard,
	"Basilisk":        SubtypeBasilisk,
	"Bat":             SubtypeBat,
	"Bear":            SubtypeBear,
	"Beast":           SubtypeBeast,
	"Beaver":          SubtypeBeaver,
	"Beeble":          SubtypeBeeble,
	"Beholder":        SubtypeBeholder,
	"Berserker":       SubtypeBerserker,
	"Bird":            SubtypeBird,
	"Blinkmoth":       SubtypeBlinkmoth,
	"Boar":            SubtypeBoar,
	"Bringer":         SubtypeBringer,
	"Brushwagg":       SubtypeBrushwagg,
	"Camarid":         SubtypeCamarid,
	"Camel":           SubtypeCamel,
	"Capybara":        SubtypeCapybara,
	"Caribou":         SubtypeCaribou,
	"Carrier":         SubtypeCarrier,
	"Cat":             SubtypeCat,
	"Centaur":         SubtypeCentaur,
	"Child":           SubtypeChild,
	"Chimera":         SubtypeChimera,
	"Citizen":         SubtypeCitizen,
	"Cleric":          SubtypeCleric,
	"Clown":           SubtypeClown,
	"Cockatrice":      SubtypeCockatrice,
	"Construct":       SubtypeConstruct,
	"Coward":          SubtypeCoward,
	"Coyote":          SubtypeCoyote,
	"Crab":            SubtypeCrab,
	"Crocodile":       SubtypeCrocodile,
	"C’tan":           SubtypeCTan,
	"Custodes":        SubtypeCustodes,
	"Cyberman":        SubtypeCyberman,
	"Cyclops":         SubtypeCyclops,
	"Dalek":           SubtypeDalek,
	"Dauthi":          SubtypeDauthi,
	"Demigod":         SubtypeDemigod,
	"Demon":           SubtypeDemon,
	"Deserter":        SubtypeDeserter,
	"Detective":       SubtypeDetective,
	"Devil":           SubtypeDevil,
	"Dinosaur":        SubtypeDinosaur,
	"Djinn":           SubtypeDjinn,
	"Doctor":          SubtypeDoctor,
	"Dog":             SubtypeDog,
	"Dragon":          SubtypeDragon,
	"Drake":           SubtypeDrake,
	"Dreadnought":     SubtypeDreadnought,
	"Drone":           SubtypeDrone,
	"Druid":           SubtypeDruid,
	"Dryad":           SubtypeDryad,
	"Dwarf":           SubtypeDwarf,
	"Efreet":          SubtypeEfreet,
	"Egg":             SubtypeEgg,
	"Elder":           SubtypeElder,
	"Eldrazi":         SubtypeEldrazi,
	"Elemental":       SubtypeElemental,
	"Elephant":        SubtypeElephant,
	"Elf":             SubtypeElf,
	"Elk":             SubtypeElk,
	"Employee":        SubtypeEmployee,
	"Eye":             SubtypeEye,
	"Faerie":          SubtypeFaerie,
	"Ferret":          SubtypeFerret,
	"Fish":            SubtypeFish,
	"Flagbearer":      SubtypeFlagbearer,
	"Fox":             SubtypeFox,
	"Fractal":         SubtypeFractal,
	"Frog":            SubtypeFrog,
	"Fungus":          SubtypeFungus,
	"Gamer":           SubtypeGamer,
	"Gargoyle":        SubtypeGargoyle,
	"Germ":            SubtypeGerm,
	"Giant":           SubtypeGiant,
	"Gith":            SubtypeGith,
	"Glimmer":         SubtypeGlimmer,
	"Gnoll":           SubtypeGnoll,
	"Gnome":           SubtypeGnome,
	"Goat":            SubtypeGoat,
	"Goblin":          SubtypeGoblin,
	"God":             SubtypeGod,
	"Golem":           SubtypeGolem,
	"Gorgon":          SubtypeGorgon,
	"Graveborn":       SubtypeGraveborn,
	"Gremlin":         SubtypeGremlin,
	"Griffin":         SubtypeGriffin,
	"Guest":           SubtypeGuest,
	"Hag":             SubtypeHag,
	"Halfling":        SubtypeHalfling,
	"Hamster":         SubtypeHamster,
	"Harpy":           SubtypeHarpy,
	"Hellion":         SubtypeHellion,
	"Hippo":           SubtypeHippo,
	"Hippogriff":      SubtypeHippogriff,
	"Homarid":         SubtypeHomarid,
	"Homunculus":      SubtypeHomunculus,
	"Horror":          SubtypeHorror,
	"Horse":           SubtypeHorse,
	"Human":           SubtypeHuman,
	"Hydra":           SubtypeHydra,
	"Hyena":           SubtypeHyena,
	"Illusion":        SubtypeIllusion,
	"Imp":             SubtypeImp,
	"Incarnation":     SubtypeIncarnation,
	"Inkling":         SubtypeInkling,
	"Inquisitor":      SubtypeInquisitor,
	"Insect":          SubtypeInsect,
	"Jackal":          SubtypeJackal,
	"Jellyfish":       SubtypeJellyfish,
	"Juggernaut":      SubtypeJuggernaut,
	"Kavu":            SubtypeKavu,
	"Kirin":           SubtypeKirin,
	"Kithkin":         SubtypeKithkin,
	"Knight":          SubtypeKnight,
	"Kobold":          SubtypeKobold,
	"Kor":             SubtypeKor,
	"Kraken":          SubtypeKraken,
	"Llama":           SubtypeLlama,
	"Lamia":           SubtypeLamia,
	"Lammasu":         SubtypeLammasu,
	"Leech":           SubtypeLeech,
	"Leviathan":       SubtypeLeviathan,
	"Lhurgoyf":        SubtypeLhurgoyf,
	"Licid":           SubtypeLicid,
	"Lizard":          SubtypeLizard,
	"Manticore":       SubtypeManticore,
	"Masticore":       SubtypeMasticore,
	"Mercenary":       SubtypeMercenary,
	"Merfolk":         SubtypeMerfolk,
	"Metathran":       SubtypeMetathran,
	"Minion":          SubtypeMinion,
	"Minotaur":        SubtypeMinotaur,
	"Mite":            SubtypeMite,
	"Mole":            SubtypeMole,
	"Monger":          SubtypeMonger,
	"Mongoose":        SubtypeMongoose,
	"Monk":            SubtypeMonk,
	"Monkey":          SubtypeMonkey,
	"Moonfolk":        SubtypeMoonfolk,
	"Mount":           SubtypeMount,
	"Mouse":           SubtypeMouse,
	"Mutant":          SubtypeMutant,
	"Myr":             SubtypeMyr,
	"Mystic":          SubtypeMystic,
	"Nautilus":        SubtypeNautilus,
	"Necron":          SubtypeNecron,
	"Nephilim":        SubtypeNephilim,
	"Nightmare":       SubtypeNightmare,
	"Nightstalker":    SubtypeNightstalker,
	"Ninja":           SubtypeNinja,
	"Noble":           SubtypeNoble,
	"Noggle":          SubtypeNoggle,
	"Nomad":           SubtypeNomad,
	"Nymph":           SubtypeNymph,
	"Octopus":         SubtypeOctopus,
	"Ogre":            SubtypeOgre,
	"Ooze":            SubtypeOoze,
	"Orb":             SubtypeOrb,
	"Orc":             SubtypeOrc,
	"Orgg":            SubtypeOrgg,
	"Otter":           SubtypeOtter,
	"Ouphe":           SubtypeOuphe,
	"Ox":              SubtypeOx,
	"Oyster":          SubtypeOyster,
	"Pangolin":        SubtypePangolin,
	"Peasant":         SubtypePeasant,
	"Pegasus":         SubtypePegasus,
	"Pentavite":       SubtypePentavite,
	"Performer":       SubtypePerformer,
	"Pest":            SubtypePest,
	"Phelddagrif":     SubtypePhelddagrif,
	"Phoenix":         SubtypePhoenix,
	"Phyrexian":       SubtypePhyrexian,
	"Pilot":           SubtypePilot,
	"Pincher":         SubtypePincher,
	"Pirate":          SubtypePirate,
	"Plant":           SubtypePlant,
	"Porcupine":       SubtypePorcupine,
	"Possum":          SubtypePossum,
	"Praetor":         SubtypePraetor,
	"Primarch":        SubtypePrimarch,
	"Prism":           SubtypePrism,
	"Processor":       SubtypeProcessor,
	"Rabbit":          SubtypeRabbit,
	"Raccoon":         SubtypeRaccoon,
	"Ranger":          SubtypeRanger,
	"Rat":             SubtypeRat,
	"Rebel":           SubtypeRebel,
	"Reflection":      SubtypeReflection,
	"Rhino":           SubtypeRhino,
	"Rigger":          SubtypeRigger,
	"Robot":           SubtypeRobot,
	"Rogue":           SubtypeRogue,
	"Sable":           SubtypeSable,
	"Salamander":      SubtypeSalamander,
	"Samurai":         SubtypeSamurai,
	"Sand":            SubtypeSand,
	"Saproling":       SubtypeSaproling,
	"Satyr":           SubtypeSatyr,
	"Scarecrow":       SubtypeScarecrow,
	"Scientist":       SubtypeScientist,
	"Scion":           SubtypeScion,
	"Scorpion":        SubtypeScorpion,
	"Scout":           SubtypeScout,
	"Sculpture":       SubtypeSculpture,
	"Serf":            SubtypeSerf,
	"Serpent":         SubtypeSerpent,
	"Servo":           SubtypeServo,
	"Shade":           SubtypeShade,
	"Shaman":          SubtypeShaman,
	"Shapeshifter":    SubtypeShapeshifter,
	"Shark":           SubtypeShark,
	"Sheep":           SubtypeSheep,
	"Siren":           SubtypeSiren,
	"Skeleton":        SubtypeSkeleton,
	"Skunk":           SubtypeSkunk,
	"Slith":           SubtypeSlith,
	"Sliver":          SubtypeSliver,
	"Sloth":           SubtypeSloth,
	"Slug":            SubtypeSlug,
	"Snail":           SubtypeSnail,
	"Snake":           SubtypeSnake,
	"Soldier":         SubtypeSoldier,
	"Soltari":         SubtypeSoltari,
	"Spawn":           SubtypeSpawn,
	"Specter":         SubtypeSpecter,
	"Spellshaper":     SubtypeSpellshaper,
	"Sphinx":          SubtypeSphinx,
	"Spider":          SubtypeSpider,
	"Spike":           SubtypeSpike,
	"Spirit":          SubtypeSpirit,
	"Splinter":        SubtypeSplinter,
	"Sponge":          SubtypeSponge,
	"Squid":           SubtypeSquid,
	"Squirrel":        SubtypeSquirrel,
	"Starfish":        SubtypeStarfish,
	"Surrakar":        SubtypeSurrakar,
	"Survivor":        SubtypeSurvivor,
	"Synth":           SubtypeSynth,
	"Tentacle":        SubtypeTentacle,
	"Tetravite":       SubtypeTetravite,
	"Thalakos":        SubtypeThalakos,
	"Thopter":         SubtypeThopter,
	"Thrull":          SubtypeThrull,
	"Tiefling":        SubtypeTiefling,
	"Toy":             SubtypeToy,
	"Treefolk":        SubtypeTreefolk,
	"Trilobite":       SubtypeTrilobite,
	"Triskelavite":    SubtypeTriskelavite,
	"Troll":           SubtypeTroll,
	"Turtle":          SubtypeTurtle,
	"Tyranid":         SubtypeTyranid,
	"Unicorn":         SubtypeUnicorn,
	"Vampire":         SubtypeVampire,
	"Varmint":         SubtypeVarmint,
	"Vedalken":        SubtypeVedalken,
	"Volver":          SubtypeVolver,
	"Wall":            SubtypeWall,
	"Walrus":          SubtypeWalrus,
	"Warlock":         SubtypeWarlock,
	"Warrior":         SubtypeWarrior,
	"Weasel":          SubtypeWeasel,
	"Weird":           SubtypeWeird,
	"Werewolf":        SubtypeWerewolf,
	"Whale":           SubtypeWhale,
	"Wizard":          SubtypeWizard,
	"Wolf":            SubtypeWolf,
	"Wolverine":       SubtypeWolverine,
	"Wombat":          SubtypeWombat,
	"Worm":            SubtypeWorm,
	"Wraith":          SubtypeWraith,
	"Wurm":            SubtypeWurm,
	"Yeti":            SubtypeYeti,
	"Zombie":          SubtypeZombie,
	"Zubera":          SubtypeZubera,
}
