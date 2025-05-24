package game

// TODO: Maybe this should be a type like CardType
// type AbilityKeyword string

const (
	AbilityKeywordAbsorb                string = "Absorb"
	AbilityKeywordAffinity              string = "Affinity"
	AbilityKeywordAfflict               string = "Afflict"
	AbilityKeywordAfterlife             string = "Afterlife"
	AbilityKeywordAftermath             string = "Aftermath"
	AbilityKeywordAmplify               string = "Amplify"
	AbilityKeywordAnnihilator           string = "Annihilator"
	AbilityKeywordAscend                string = "Ascend"
	AbilityKeywordAssist                string = "Assist"
	AbilityKeywordAuraSwap              string = "Aura Swap"
	AbilityKeywordAwaken                string = "Awaken"
	AbilityKeywordBackup                string = "Backup"
	AbilityKeywordBanding               string = "Banding"
	AbilityKeywordBargain               string = "Bargain"
	AbilityKeywordBattleCry             string = "Battle Cry"
	AbilityKeywordBestow                string = "Bestow"
	AbilityKeywordBlitz                 string = "Blitz"
	AbilityKeywordBloodthirst           string = "Bloodthirst"
	AbilityKeywordBoast                 string = "Boast"
	AbilityKeywordBushido               string = "Bushido"
	AbilityKeywordBuyback               string = "Buyback"
	AbilityKeywordCascade               string = "Cascade"
	AbilityKeywordCasualty              string = "Casualty"
	AbilityKeywordChampion              string = "Champion"
	AbilityKeywordChangeling            string = "Changeling"
	AbilityKeywordCipher                string = "Cipher"
	AbilityKeywordCleave                string = "Cleave"
	AbilityKeywordCompanion             string = "Companion"
	AbilityKeywordCompleated            string = "Compleated"
	AbilityKeywordConspire              string = "Conspire"
	AbilityKeywordConvoke               string = "Convoke"
	AbilityKeywordCraft                 string = "Craft"
	AbilityKeywordCrew                  string = "Crew"
	AbilityKeywordCumulativeUpkeep      string = "Cumulative Upkeep"
	AbilityKeywordCycling               string = "Cycling"
	AbilityKeywordDash                  string = "Dash"
	AbilityKeywordDayboundAndNightbound string = "Daybound and Nightbound"
	AbilityKeywordDeathtouch            string = "Deathtouch"
	AbilityKeywordDecayed               string = "Decayed"
	AbilityKeywordDefender              string = "Defender"
	AbilityKeywordDelve                 string = "Delve"
	AbilityKeywordDemonstrate           string = "Demonstrate"
	AbilityKeywordDethrone              string = "Dethrone"
	AbilityKeywordDevoid                string = "Devoid"
	AbilityKeywordDevour                string = "Devour"
	AbilityKeywordDisguise              string = "Disguise"
	AbilityKeywordDisturb               string = "Disturb"
	AbilityKeywordDoubleStrike          string = "Double Strike"
	AbilityKeywordDredge                string = "Dredge"
	AbilityKeywordEcho                  string = "Echo"
	AbilityKeywordEmbalm                string = "Embalm"
	AbilityKeywordEmerge                string = "Emerge"
	AbilityKeywordEnchant               string = "Enchant"
	AbilityKeywordEncore                string = "Encore"
	AbilityKeywordEnlist                string = "Enlist"
	AbilityKeywordEntwine               string = "Entwine"
	AbilityKeywordEpic                  string = "Epic"
	AbilityKeywordEquip                 string = "Equip"
	AbilityKeywordEscalate              string = "Escalate"
	AbilityKeywordEscape                string = "Escape"
	AbilityKeywordEternalize            string = "Eternalize"
	AbilityKeywordEvoke                 string = "Evoke"
	AbilityKeywordEvolve                string = "Evolve"
	AbilityKeywordExalted               string = "Exalted"
	AbilityKeywordExploit               string = "Exploit"
	AbilityKeywordExtort                string = "Extort"
	AbilityKeywordFabricate             string = "Fabricate"
	AbilityKeywordFading                string = "Fading"
	AbilityKeywordFear                  string = "Fear"
	AbilityKeywordFirstStrike           string = "First Strike"
	AbilityKeywordFlanking              string = "Flanking"
	AbilityKeywordFlash                 string = "Flash"
	AbilityKeywordFlashback             string = "Flashback"
	AbilityKeywordFlying                string = "Flying"
	AbilityKeywordForMirrodin           string = "For Mirrodin!"
	AbilityKeywordForecast              string = "Forecast"
	AbilityKeywordForetell              string = "Foretell"
	AbilityKeywordFortify               string = "Fortify"
	AbilityKeywordFreerunning           string = "Freerunning"
	AbilityKeywordFrenzy                string = "Frenzy"
	AbilityKeywordFuse                  string = "Fuse"
	AbilityKeywordGift                  string = "Gift"
	AbilityKeywordGraft                 string = "Graft"
	AbilityKeywordGravestorm            string = "Gravestorm"
	AbilityKeywordHaste                 string = "Haste"
	AbilityKeywordHaunt                 string = "Haunt"
	AbilityKeywordHexproof              string = "Hexproof"
	AbilityKeywordHiddenAgenda          string = "Hidden Agenda"
	AbilityKeywordHideaway              string = "Hideaway"
	AbilityKeywordHorsemanship          string = "Horsemanship"
	AbilityKeywordImpending             string = "Impending"
	AbilityKeywordImprovise             string = "Improvise"
	AbilityKeywordIndestructible        string = "Indestructible"
	AbilityKeywordInfect                string = "Infect"
	AbilityKeywordIngest                string = "Ingest"
	AbilityKeywordIntimidate            string = "Intimidate"
	AbilityKeywordJumpStart             string = "Jump-Start"
	AbilityKeywordKicker                string = "Kicker"
	AbilityKeywordLandwalk              string = "Landwalk"
	AbilityKeywordLevelUp               string = "Level Up"
	AbilityKeywordLifelink              string = "Lifelink"
	AbilityKeywordLivingMetal           string = "Living Metal"
	AbilityKeywordLivingWeapon          string = "Living Weapon"
	AbilityKeywordMadness               string = "Madness"
	AbilityKeywordMelee                 string = "Melee"
	AbilityKeywordMenace                string = "Menace"
	AbilityKeywordMentor                string = "Mentor"
	AbilityKeywordMiracle               string = "Miracle"
	AbilityKeywordModular               string = "Modular"
	AbilityKeywordMoreThanMeetsTheEye   string = "More Than Meets the Eye"
	AbilityKeywordMorph                 string = "Morph"
	AbilityKeywordMutate                string = "Mutate"
	AbilityKeywordMyriad                string = "Myriad"
	AbilityKeywordNinjutsu              string = "Ninjutsu"
	AbilityKeywordOffering              string = "Offering"
	AbilityKeywordOffspring             string = "Offspring"
	AbilityKeywordOutlast               string = "Outlast"
	AbilityKeywordOverload              string = "Overload"
	AbilityKeywordPartner               string = "Partner"
	AbilityKeywordPersist               string = "Persist"
	AbilityKeywordPhasing               string = "Phasing"
	AbilityKeywordPlot                  string = "Plot"
	AbilityKeywordPoisonous             string = "Poisonous"
	AbilityKeywordProtection            string = "Protection"
	AbilityKeywordPrototype             string = "Prototype"
	AbilityKeywordProvoke               string = "Provoke"
	AbilityKeywordProwess               string = "Prowess"
	AbilityKeywordProwl                 string = "Prowl"
	AbilityKeywordRampage               string = "Rampage"
	AbilityKeywordRavenous              string = "Ravenous"
	AbilityKeywordReach                 string = "Reach"
	AbilityKeywordReadAhead             string = "Read Ahead"
	AbilityKeywordRebound               string = "Rebound"
	AbilityKeywordReconfigure           string = "Reconfigure"
	AbilityKeywordRecover               string = "Recover"
	AbilityKeywordReinforce             string = "Reinforce"
	AbilityKeywordRenown                string = "Renown"
	AbilityKeywordReplicate             string = "Replicate"
	AbilityKeywordRetrace               string = "Retrace"
	AbilityKeywordRiot                  string = "Riot"
	AbilityKeywordRipple                string = "Ripple"
	AbilityKeywordSaddle                string = "Saddle"
	AbilityKeywordScavenge              string = "Scavenge"
	AbilityKeywordShadow                string = "Shadow"
	AbilityKeywordShroud                string = "Shroud"
	AbilityKeywordSkulk                 string = "Skulk"
	AbilityKeywordSolved                string = "Solved"
	AbilityKeywordSoulbond              string = "Soulbond"
	AbilityKeywordSoulshift             string = "Soulshift"
	AbilityKeywordSpaceSculptor         string = "Space Sculptor"
	AbilityKeywordSpectacle             string = "Spectacle"
	AbilityKeywordSplice                string = "Splice"
	AbilityKeywordSplitSecond           string = "Split Second"
	AbilityKeywordSpree                 string = "Spree"
	AbilityKeywordSquad                 string = "Squad"
	AbilityKeywordStorm                 string = "Storm"
	AbilityKeywordSunburst              string = "Sunburst"
	AbilityKeywordSurge                 string = "Surge"
	AbilityKeywordSuspend               string = "Suspend"
	AbilityKeywordToxic                 string = "Toxic"
	AbilityKeywordTraining              string = "Training"
	AbilityKeywordTrample               string = "Trample"
	AbilityKeywordTransfigure           string = "Transfigure"
	AbilityKeywordTransmute             string = "Transmute"
	AbilityKeywordTribute               string = "Tribute"
	AbilityKeywordUmbraArmor            string = "Umbra Armor"
	AbilityKeywordUndaunted             string = "Undaunted"
	AbilityKeywordUndying               string = "Undying"
	AbilityKeywordUnearth               string = "Unearth"
	AbilityKeywordUnleash               string = "Unleash"
	AbilityKeywordVanishing             string = "Vanishing"
	AbilityKeywordVigilance             string = "Vigilance"
	AbilityKeywordVisit                 string = "Visit"
	AbilityKeywordWard                  string = "Ward"
	AbilityKeywordWither                string = "Wither"
)

var KeywordCombatAbilities = []string{
	AbilityKeywordDeathtouch,
	AbilityKeywordDoubleStrike,
	AbilityKeywordFirstStrike,
	AbilityKeywordFlying,
	AbilityKeywordHaste,
	AbilityKeywordLifelink,
	AbilityKeywordReach,
	AbilityKeywordTrample,
	AbilityKeywordVigilance,
	AbilityKeywordMenace,
	AbilityKeywordProtection,
	AbilityKeywordHexproof,
	AbilityKeywordShroud,
	AbilityKeywordIndestructible,
	AbilityKeywordFear,
	AbilityKeywordIntimidate,
	AbilityKeywordSkulk,
}

var KeywordCostModifiers = []string{
	AbilityKeywordAffinity,
	AbilityKeywordConvoke,
	AbilityKeywordDelve,
	AbilityKeywordImprovise,
	AbilityKeywordUndaunted,
	AbilityKeywordEscape,
}

var KeywordAlternativeCosts = []string{
	AbilityKeywordBuyback,
	AbilityKeywordFlashback,
	AbilityKeywordMadness,
	AbilityKeywordDash,
	AbilityKeywordBestow,
	AbilityKeywordOverload,
	AbilityKeywordFuse,
	AbilityKeywordJumpStart,
	AbilityKeywordForetell,
	AbilityKeywordReplicate,
	AbilityKeywordRetrace,
	AbilityKeywordAftermath,
	AbilityKeywordDemonstrate,
}

var KeywordEvasion = []string{
	AbilityKeywordFlying,
	AbilityKeywordHorsemanship,
	AbilityKeywordShadow,
	AbilityKeywordLandwalk,
	AbilityKeywordFear,
	AbilityKeywordIntimidate,
	AbilityKeywordSkulk,
}

var KeywordTriggeredAbilities = []string{
	AbilityKeywordCascade,
	AbilityKeywordDredge,
	AbilityKeywordPersist,
	AbilityKeywordUndying,
	AbilityKeywordUnearth,
	AbilityKeywordRebound,
	AbilityKeywordRetrace,
	AbilityKeywordMiracle,
	AbilityKeywordExploit,
	AbilityKeywordRenown,
	AbilityKeywordBoast,
	AbilityKeywordEvoke,
}

var KeywordStaticAbilities = []string{
	AbilityKeywordDefender,
	AbilityKeywordChangeling,
	AbilityKeywordDevoid,
	AbilityKeywordDayboundAndNightbound,
	AbilityKeywordIndestructible,
	AbilityKeywordShroud,
	AbilityKeywordHexproof,
	AbilityKeywordReach,
	AbilityKeywordMenace,
}

func IsCombatAbilities(kw string) bool {
	for _, k := range KeywordCombatAbilities {
		if k == kw {
			return true
		}
	}
	return false
}

func IsCostModifiers(kw string) bool {
	for _, k := range KeywordCostModifiers {
		if k == kw {
			return true
		}
	}
	return false
}

func IsAlternativeCosts(kw string) bool {
	for _, k := range KeywordAlternativeCosts {
		if k == kw {
			return true
		}
	}
	return false
}

func IsEvasion(kw string) bool {
	for _, k := range KeywordEvasion {
		if k == kw {
			return true
		}
	}
	return false
}

func IsTriggeredAbilities(kw string) bool {
	for _, k := range KeywordTriggeredAbilities {
		if k == kw {
			return true
		}
	}
	return false
}

func IsStaticAbilities(kw string) bool {
	for _, k := range KeywordStaticAbilities {
		if k == kw {
			return true
		}
	}
	return false
}

var AllAbilityKeywords = map[string]string{
	"Deathtouch":              AbilityKeywordDeathtouch,
	"Defender":                AbilityKeywordDefender,
	"Double Strike":           AbilityKeywordDoubleStrike,
	"Enchant":                 AbilityKeywordEnchant,
	"Equip":                   AbilityKeywordEquip,
	"First Strike":            AbilityKeywordFirstStrike,
	"Flash":                   AbilityKeywordFlash,
	"Flying":                  AbilityKeywordFlying,
	"Haste":                   AbilityKeywordHaste,
	"Hexproof":                AbilityKeywordHexproof,
	"Indestructible":          AbilityKeywordIndestructible,
	"Intimidate":              AbilityKeywordIntimidate,
	"Landwalk":                AbilityKeywordLandwalk,
	"Lifelink":                AbilityKeywordLifelink,
	"Protection":              AbilityKeywordProtection,
	"Reach":                   AbilityKeywordReach,
	"Shroud":                  AbilityKeywordShroud,
	"Trample":                 AbilityKeywordTrample,
	"Vigilance":               AbilityKeywordVigilance,
	"Ward":                    AbilityKeywordWard,
	"Banding":                 AbilityKeywordBanding,
	"Rampage":                 AbilityKeywordRampage,
	"Cumulative Upkeep":       AbilityKeywordCumulativeUpkeep,
	"Flanking":                AbilityKeywordFlanking,
	"Phasing":                 AbilityKeywordPhasing,
	"Buyback":                 AbilityKeywordBuyback,
	"Shadow":                  AbilityKeywordShadow,
	"Cycling":                 AbilityKeywordCycling,
	"Echo":                    AbilityKeywordEcho,
	"Horsemanship":            AbilityKeywordHorsemanship,
	"Fading":                  AbilityKeywordFading,
	"Kicker":                  AbilityKeywordKicker,
	"Flashback":               AbilityKeywordFlashback,
	"Madness":                 AbilityKeywordMadness,
	"Fear":                    AbilityKeywordFear,
	"Morph":                   AbilityKeywordMorph,
	"Amplify":                 AbilityKeywordAmplify,
	"Provoke":                 AbilityKeywordProvoke,
	"Storm":                   AbilityKeywordStorm,
	"Affinity":                AbilityKeywordAffinity,
	"Entwine":                 AbilityKeywordEntwine,
	"Modular":                 AbilityKeywordModular,
	"Sunburst":                AbilityKeywordSunburst,
	"Bushido":                 AbilityKeywordBushido,
	"Soulshift":               AbilityKeywordSoulshift,
	"Splice":                  AbilityKeywordSplice,
	"Offering":                AbilityKeywordOffering,
	"Ninjutsu":                AbilityKeywordNinjutsu,
	"Epic":                    AbilityKeywordEpic,
	"Convoke":                 AbilityKeywordConvoke,
	"Dredge":                  AbilityKeywordDredge,
	"Transmute":               AbilityKeywordTransmute,
	"Bloodthirst":             AbilityKeywordBloodthirst,
	"Haunt":                   AbilityKeywordHaunt,
	"Replicate":               AbilityKeywordReplicate,
	"Forecast":                AbilityKeywordForecast,
	"Graft":                   AbilityKeywordGraft,
	"Recover":                 AbilityKeywordRecover,
	"Ripple":                  AbilityKeywordRipple,
	"Split Second":            AbilityKeywordSplitSecond,
	"Suspend":                 AbilityKeywordSuspend,
	"Vanishing":               AbilityKeywordVanishing,
	"Absorb":                  AbilityKeywordAbsorb,
	"Aura Swap":               AbilityKeywordAuraSwap,
	"Delve":                   AbilityKeywordDelve,
	"Fortify":                 AbilityKeywordFortify,
	"Frenzy":                  AbilityKeywordFrenzy,
	"Gravestorm":              AbilityKeywordGravestorm,
	"Poisonous":               AbilityKeywordPoisonous,
	"Transfigure":             AbilityKeywordTransfigure,
	"Champion":                AbilityKeywordChampion,
	"Changeling":              AbilityKeywordChangeling,
	"Evoke":                   AbilityKeywordEvoke,
	"Hideaway":                AbilityKeywordHideaway,
	"Prowl":                   AbilityKeywordProwl,
	"Reinforce":               AbilityKeywordReinforce,
	"Conspire":                AbilityKeywordConspire,
	"Persist":                 AbilityKeywordPersist,
	"Wither":                  AbilityKeywordWither,
	"Retrace":                 AbilityKeywordRetrace,
	"Devour":                  AbilityKeywordDevour,
	"Exalted":                 AbilityKeywordExalted,
	"Unearth":                 AbilityKeywordUnearth,
	"Cascade":                 AbilityKeywordCascade,
	"Annihilator":             AbilityKeywordAnnihilator,
	"Level Up":                AbilityKeywordLevelUp,
	"Rebound":                 AbilityKeywordRebound,
	"Umbra Armor":             AbilityKeywordUmbraArmor,
	"Infect":                  AbilityKeywordInfect,
	"Battle Cry":              AbilityKeywordBattleCry,
	"Living Weapon":           AbilityKeywordLivingWeapon,
	"Undying":                 AbilityKeywordUndying,
	"Miracle":                 AbilityKeywordMiracle,
	"Soulbond":                AbilityKeywordSoulbond,
	"Overload":                AbilityKeywordOverload,
	"Scavenge":                AbilityKeywordScavenge,
	"Unleash":                 AbilityKeywordUnleash,
	"Cipher":                  AbilityKeywordCipher,
	"Evolve":                  AbilityKeywordEvolve,
	"Extort":                  AbilityKeywordExtort,
	"Fuse":                    AbilityKeywordFuse,
	"Bestow":                  AbilityKeywordBestow,
	"Tribute":                 AbilityKeywordTribute,
	"Dethrone":                AbilityKeywordDethrone,
	"Hidden Agenda":           AbilityKeywordHiddenAgenda,
	"Outlast":                 AbilityKeywordOutlast,
	"Prowess":                 AbilityKeywordProwess,
	"Dash":                    AbilityKeywordDash,
	"Exploit":                 AbilityKeywordExploit,
	"Menace":                  AbilityKeywordMenace,
	"Renown":                  AbilityKeywordRenown,
	"Awaken":                  AbilityKeywordAwaken,
	"Devoid":                  AbilityKeywordDevoid,
	"Ingest":                  AbilityKeywordIngest,
	"Myriad":                  AbilityKeywordMyriad,
	"Surge":                   AbilityKeywordSurge,
	"Skulk":                   AbilityKeywordSkulk,
	"Emerge":                  AbilityKeywordEmerge,
	"Escalate":                AbilityKeywordEscalate,
	"Melee":                   AbilityKeywordMelee,
	"Crew":                    AbilityKeywordCrew,
	"Fabricate":               AbilityKeywordFabricate,
	"Partner":                 AbilityKeywordPartner,
	"Undaunted":               AbilityKeywordUndaunted,
	"Improvise":               AbilityKeywordImprovise,
	"Aftermath":               AbilityKeywordAftermath,
	"Embalm":                  AbilityKeywordEmbalm,
	"Eternalize":              AbilityKeywordEternalize,
	"Afflict":                 AbilityKeywordAfflict,
	"Ascend":                  AbilityKeywordAscend,
	"Assist":                  AbilityKeywordAssist,
	"Jump-Start":              AbilityKeywordJumpStart,
	"Mentor":                  AbilityKeywordMentor,
	"Afterlife":               AbilityKeywordAfterlife,
	"Riot":                    AbilityKeywordRiot,
	"Spectacle":               AbilityKeywordSpectacle,
	"Escape":                  AbilityKeywordEscape,
	"Companion":               AbilityKeywordCompanion,
	"Mutate":                  AbilityKeywordMutate,
	"Encore":                  AbilityKeywordEncore,
	"Boast":                   AbilityKeywordBoast,
	"Foretell":                AbilityKeywordForetell,
	"Demonstrate":             AbilityKeywordDemonstrate,
	"Daybound and Nightbound": AbilityKeywordDayboundAndNightbound,
	"Disturb":                 AbilityKeywordDisturb,
	"Decayed":                 AbilityKeywordDecayed,
	"Cleave":                  AbilityKeywordCleave,
	"Training":                AbilityKeywordTraining,
	"Compleated":              AbilityKeywordCompleated,
	"Reconfigure":             AbilityKeywordReconfigure,
	"Blitz":                   AbilityKeywordBlitz,
	"Casualty":                AbilityKeywordCasualty,
	"Enlist":                  AbilityKeywordEnlist,
	"Read Ahead":              AbilityKeywordReadAhead,
	"Ravenous":                AbilityKeywordRavenous,
	"Squad":                   AbilityKeywordSquad,
	"Space Sculptor":          AbilityKeywordSpaceSculptor,
	"Visit":                   AbilityKeywordVisit,
	"Prototype":               AbilityKeywordPrototype,
	"Living Metal":            AbilityKeywordLivingMetal,
	"More Than Meets the Eye": AbilityKeywordMoreThanMeetsTheEye,
	"For Mirrodin!":           AbilityKeywordForMirrodin,
	"Toxic":                   AbilityKeywordToxic,
	"Backup":                  AbilityKeywordBackup,
	"Bargain":                 AbilityKeywordBargain,
	"Craft":                   AbilityKeywordCraft,
	"Disguise":                AbilityKeywordDisguise,
	"Solved":                  AbilityKeywordSolved,
	"Plot":                    AbilityKeywordPlot,
	"Saddle":                  AbilityKeywordSaddle,
	"Spree":                   AbilityKeywordSpree,
	"Freerunning":             AbilityKeywordFreerunning,
	"Gift":                    AbilityKeywordGift,
	"Offspring":               AbilityKeywordOffspring,
	"Impending":               AbilityKeywordImpending,
}

func IsAbilityKeyword(s string) bool {
	_, ok := AllAbilityKeywords[s]
	return ok
}

func ParseAbilityKeyword(s string) (string, bool) {
	kw, ok := AllAbilityKeywords[s]
	return kw, ok
}
