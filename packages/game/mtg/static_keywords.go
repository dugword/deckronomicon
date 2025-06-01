package mtg

type StaticKeyword string

// TODO Maybe have this just be "StaticKeywordAbsorb"
const (
	StaticKeywordAbsorb                StaticKeyword = "Absorb"
	StaticKeywordAffinity              StaticKeyword = "Affinity"
	StaticKeywordAfflict               StaticKeyword = "Afflict"
	StaticKeywordAfterlife             StaticKeyword = "Afterlife"
	StaticKeywordAftermath             StaticKeyword = "Aftermath"
	StaticKeywordAmplify               StaticKeyword = "Amplify"
	StaticKeywordAnnihilator           StaticKeyword = "Annihilator"
	StaticKeywordAscend                StaticKeyword = "Ascend"
	StaticKeywordAssist                StaticKeyword = "Assist"
	StaticKeywordAuraSwap              StaticKeyword = "Aura Swap"
	StaticKeywordAwaken                StaticKeyword = "Awaken"
	StaticKeywordBackup                StaticKeyword = "Backup"
	StaticKeywordBanding               StaticKeyword = "Banding"
	StaticKeywordBargain               StaticKeyword = "Bargain"
	StaticKeywordBattleCry             StaticKeyword = "Battle Cry"
	StaticKeywordBestow                StaticKeyword = "Bestow"
	StaticKeywordBlitz                 StaticKeyword = "Blitz"
	StaticKeywordBloodthirst           StaticKeyword = "Bloodthirst"
	StaticKeywordBoast                 StaticKeyword = "Boast"
	StaticKeywordBushido               StaticKeyword = "Bushido"
	StaticKeywordBuyback               StaticKeyword = "Buyback"
	StaticKeywordCascade               StaticKeyword = "Cascade"
	StaticKeywordCasualty              StaticKeyword = "Casualty"
	StaticKeywordChampion              StaticKeyword = "Champion"
	StaticKeywordChangeling            StaticKeyword = "Changeling"
	StaticKeywordCipher                StaticKeyword = "Cipher"
	StaticKeywordCleave                StaticKeyword = "Cleave"
	StaticKeywordCompanion             StaticKeyword = "Companion"
	StaticKeywordCompleated            StaticKeyword = "Compleated"
	StaticKeywordConspire              StaticKeyword = "Conspire"
	StaticKeywordConvoke               StaticKeyword = "Convoke"
	StaticKeywordCraft                 StaticKeyword = "Craft"
	StaticKeywordCrew                  StaticKeyword = "Crew"
	StaticKeywordCumulativeUpkeep      StaticKeyword = "Cumulative Upkeep"
	StaticKeywordCycling               StaticKeyword = "Cycling"
	StaticKeywordDash                  StaticKeyword = "Dash"
	StaticKeywordDayboundAndNightbound StaticKeyword = "Daybound and Nightbound"
	StaticKeywordDeathtouch            StaticKeyword = "Deathtouch"
	StaticKeywordDecayed               StaticKeyword = "Decayed"
	StaticKeywordDefender              StaticKeyword = "Defender"
	StaticKeywordDelve                 StaticKeyword = "Delve"
	StaticKeywordDemonstrate           StaticKeyword = "Demonstrate"
	StaticKeywordDethrone              StaticKeyword = "Dethrone"
	StaticKeywordDevoid                StaticKeyword = "Devoid"
	StaticKeywordDevour                StaticKeyword = "Devour"
	StaticKeywordDisguise              StaticKeyword = "Disguise"
	StaticKeywordDisturb               StaticKeyword = "Disturb"
	StaticKeywordDoubleStrike          StaticKeyword = "Double Strike"
	StaticKeywordDredge                StaticKeyword = "Dredge"
	StaticKeywordEcho                  StaticKeyword = "Echo"
	StaticKeywordEmbalm                StaticKeyword = "Embalm"
	StaticKeywordEmerge                StaticKeyword = "Emerge"
	StaticKeywordEnchant               StaticKeyword = "Enchant"
	StaticKeywordEncore                StaticKeyword = "Encore"
	StaticKeywordEnlist                StaticKeyword = "Enlist"
	StaticKeywordEntwine               StaticKeyword = "Entwine"
	StaticKeywordEpic                  StaticKeyword = "Epic"
	StaticKeywordEquip                 StaticKeyword = "Equip"
	StaticKeywordEscalate              StaticKeyword = "Escalate"
	StaticKeywordEscape                StaticKeyword = "Escape"
	StaticKeywordEternalize            StaticKeyword = "Eternalize"
	StaticKeywordEvoke                 StaticKeyword = "Evoke"
	StaticKeywordEvolve                StaticKeyword = "Evolve"
	StaticKeywordExalted               StaticKeyword = "Exalted"
	StaticKeywordExploit               StaticKeyword = "Exploit"
	StaticKeywordExtort                StaticKeyword = "Extort"
	StaticKeywordFabricate             StaticKeyword = "Fabricate"
	StaticKeywordFading                StaticKeyword = "Fading"
	StaticKeywordFear                  StaticKeyword = "Fear"
	StaticKeywordFirstStrike           StaticKeyword = "First Strike"
	StaticKeywordFlanking              StaticKeyword = "Flanking"
	StaticKeywordFlash                 StaticKeyword = "Flash"
	StaticKeywordFlashback             StaticKeyword = "Flashback"
	StaticKeywordFlying                StaticKeyword = "Flying"
	StaticKeywordForMirrodin           StaticKeyword = "For Mirrodin!"
	StaticKeywordForecast              StaticKeyword = "Forecast"
	StaticKeywordForetell              StaticKeyword = "Foretell"
	StaticKeywordFortify               StaticKeyword = "Fortify"
	StaticKeywordFreerunning           StaticKeyword = "Freerunning"
	StaticKeywordFrenzy                StaticKeyword = "Frenzy"
	StaticKeywordFuse                  StaticKeyword = "Fuse"
	StaticKeywordGift                  StaticKeyword = "Gift"
	StaticKeywordGraft                 StaticKeyword = "Graft"
	StaticKeywordGravestorm            StaticKeyword = "Gravestorm"
	StaticKeywordHaste                 StaticKeyword = "Haste"
	StaticKeywordHaunt                 StaticKeyword = "Haunt"
	StaticKeywordHexproof              StaticKeyword = "Hexproof"
	StaticKeywordHiddenAgenda          StaticKeyword = "Hidden Agenda"
	StaticKeywordHideaway              StaticKeyword = "Hideaway"
	StaticKeywordHorsemanship          StaticKeyword = "Horsemanship"
	StaticKeywordImpending             StaticKeyword = "Impending"
	StaticKeywordImprovise             StaticKeyword = "Improvise"
	StaticKeywordIndestructible        StaticKeyword = "Indestructible"
	StaticKeywordInfect                StaticKeyword = "Infect"
	StaticKeywordIngest                StaticKeyword = "Ingest"
	StaticKeywordIntimidate            StaticKeyword = "Intimidate"
	StaticKeywordJumpStart             StaticKeyword = "Jump-Start"
	StaticKeywordKicker                StaticKeyword = "Kicker"
	StaticKeywordLandwalk              StaticKeyword = "Landwalk"
	StaticKeywordLevelUp               StaticKeyword = "Level Up"
	StaticKeywordLifelink              StaticKeyword = "Lifelink"
	StaticKeywordLivingMetal           StaticKeyword = "Living Metal"
	StaticKeywordLivingWeapon          StaticKeyword = "Living Weapon"
	StaticKeywordMadness               StaticKeyword = "Madness"
	StaticKeywordMelee                 StaticKeyword = "Melee"
	StaticKeywordMenace                StaticKeyword = "Menace"
	StaticKeywordMentor                StaticKeyword = "Mentor"
	StaticKeywordMiracle               StaticKeyword = "Miracle"
	StaticKeywordModular               StaticKeyword = "Modular"
	StaticKeywordMoreThanMeetsTheEye   StaticKeyword = "More Than Meets the Eye"
	StaticKeywordMorph                 StaticKeyword = "Morph"
	StaticKeywordMutate                StaticKeyword = "Mutate"
	StaticKeywordMyriad                StaticKeyword = "Myriad"
	StaticKeywordNinjutsu              StaticKeyword = "Ninjutsu"
	StaticKeywordOffering              StaticKeyword = "Offering"
	StaticKeywordOffspring             StaticKeyword = "Offspring"
	StaticKeywordOutlast               StaticKeyword = "Outlast"
	StaticKeywordOverload              StaticKeyword = "Overload"
	StaticKeywordPartner               StaticKeyword = "Partner"
	StaticKeywordPersist               StaticKeyword = "Persist"
	StaticKeywordPhasing               StaticKeyword = "Phasing"
	StaticKeywordPlot                  StaticKeyword = "Plot"
	StaticKeywordPoisonous             StaticKeyword = "Poisonous"
	StaticKeywordProtection            StaticKeyword = "Protection"
	StaticKeywordPrototype             StaticKeyword = "Prototype"
	StaticKeywordProvoke               StaticKeyword = "Provoke"
	StaticKeywordProwess               StaticKeyword = "Prowess"
	StaticKeywordProwl                 StaticKeyword = "Prowl"
	StaticKeywordRampage               StaticKeyword = "Rampage"
	StaticKeywordRavenous              StaticKeyword = "Ravenous"
	StaticKeywordReach                 StaticKeyword = "Reach"
	StaticKeywordReadAhead             StaticKeyword = "Read Ahead"
	StaticKeywordRebound               StaticKeyword = "Rebound"
	StaticKeywordReconfigure           StaticKeyword = "Reconfigure"
	StaticKeywordRecover               StaticKeyword = "Recover"
	StaticKeywordReinforce             StaticKeyword = "Reinforce"
	StaticKeywordRenown                StaticKeyword = "Renown"
	StaticKeywordReplicate             StaticKeyword = "Replicate"
	StaticKeywordRetrace               StaticKeyword = "Retrace"
	StaticKeywordRiot                  StaticKeyword = "Riot"
	StaticKeywordRipple                StaticKeyword = "Ripple"
	StaticKeywordSaddle                StaticKeyword = "Saddle"
	StaticKeywordScavenge              StaticKeyword = "Scavenge"
	StaticKeywordShadow                StaticKeyword = "Shadow"
	StaticKeywordShroud                StaticKeyword = "Shroud"
	StaticKeywordSkulk                 StaticKeyword = "Skulk"
	StaticKeywordSolved                StaticKeyword = "Solved"
	StaticKeywordSoulbond              StaticKeyword = "Soulbond"
	StaticKeywordSoulshift             StaticKeyword = "Soulshift"
	StaticKeywordSpaceSculptor         StaticKeyword = "Space Sculptor"
	StaticKeywordSpectacle             StaticKeyword = "Spectacle"
	StaticKeywordSplice                StaticKeyword = "Splice"
	StaticKeywordSplitSecond           StaticKeyword = "Split Second"
	StaticKeywordSpree                 StaticKeyword = "Spree"
	StaticKeywordSquad                 StaticKeyword = "Squad"
	StaticKeywordStorm                 StaticKeyword = "Storm"
	StaticKeywordSunburst              StaticKeyword = "Sunburst"
	StaticKeywordSurge                 StaticKeyword = "Surge"
	StaticKeywordSuspend               StaticKeyword = "Suspend"
	StaticKeywordToxic                 StaticKeyword = "Toxic"
	StaticKeywordTraining              StaticKeyword = "Training"
	StaticKeywordTrample               StaticKeyword = "Trample"
	StaticKeywordTransfigure           StaticKeyword = "Transfigure"
	StaticKeywordTransmute             StaticKeyword = "Transmute"
	StaticKeywordTribute               StaticKeyword = "Tribute"
	StaticKeywordUmbraArmor            StaticKeyword = "Umbra Armor"
	StaticKeywordUndaunted             StaticKeyword = "Undaunted"
	StaticKeywordUndying               StaticKeyword = "Undying"
	StaticKeywordUnearth               StaticKeyword = "Unearth"
	StaticKeywordUnleash               StaticKeyword = "Unleash"
	StaticKeywordVanishing             StaticKeyword = "Vanishing"
	StaticKeywordVigilance             StaticKeyword = "Vigilance"
	StaticKeywordVisit                 StaticKeyword = "Visit"
	StaticKeywordWard                  StaticKeyword = "Ward"
	StaticKeywordWither                StaticKeyword = "Wither"
)

func StringToStaticKeyword(s string) (StaticKeyword, bool) {
	allStaticKeywords := map[string]StaticKeyword{
		"Deathtouch":              StaticKeywordDeathtouch,
		"Defender":                StaticKeywordDefender,
		"Double Strike":           StaticKeywordDoubleStrike,
		"Enchant":                 StaticKeywordEnchant,
		"Equip":                   StaticKeywordEquip,
		"First Strike":            StaticKeywordFirstStrike,
		"Flash":                   StaticKeywordFlash,
		"Flying":                  StaticKeywordFlying,
		"Haste":                   StaticKeywordHaste,
		"Hexproof":                StaticKeywordHexproof,
		"Indestructible":          StaticKeywordIndestructible,
		"Intimidate":              StaticKeywordIntimidate,
		"Landwalk":                StaticKeywordLandwalk,
		"Lifelink":                StaticKeywordLifelink,
		"Protection":              StaticKeywordProtection,
		"Reach":                   StaticKeywordReach,
		"Shroud":                  StaticKeywordShroud,
		"Trample":                 StaticKeywordTrample,
		"Vigilance":               StaticKeywordVigilance,
		"Ward":                    StaticKeywordWard,
		"Banding":                 StaticKeywordBanding,
		"Rampage":                 StaticKeywordRampage,
		"Cumulative Upkeep":       StaticKeywordCumulativeUpkeep,
		"Flanking":                StaticKeywordFlanking,
		"Phasing":                 StaticKeywordPhasing,
		"Buyback":                 StaticKeywordBuyback,
		"Shadow":                  StaticKeywordShadow,
		"Cycling":                 StaticKeywordCycling,
		"Echo":                    StaticKeywordEcho,
		"Horsemanship":            StaticKeywordHorsemanship,
		"Fading":                  StaticKeywordFading,
		"Kicker":                  StaticKeywordKicker,
		"Flashback":               StaticKeywordFlashback,
		"Madness":                 StaticKeywordMadness,
		"Fear":                    StaticKeywordFear,
		"Morph":                   StaticKeywordMorph,
		"Amplify":                 StaticKeywordAmplify,
		"Provoke":                 StaticKeywordProvoke,
		"Storm":                   StaticKeywordStorm,
		"Affinity":                StaticKeywordAffinity,
		"Entwine":                 StaticKeywordEntwine,
		"Modular":                 StaticKeywordModular,
		"Sunburst":                StaticKeywordSunburst,
		"Bushido":                 StaticKeywordBushido,
		"Soulshift":               StaticKeywordSoulshift,
		"Splice":                  StaticKeywordSplice,
		"Offering":                StaticKeywordOffering,
		"Ninjutsu":                StaticKeywordNinjutsu,
		"Epic":                    StaticKeywordEpic,
		"Convoke":                 StaticKeywordConvoke,
		"Dredge":                  StaticKeywordDredge,
		"Transmute":               StaticKeywordTransmute,
		"Bloodthirst":             StaticKeywordBloodthirst,
		"Haunt":                   StaticKeywordHaunt,
		"Replicate":               StaticKeywordReplicate,
		"Forecast":                StaticKeywordForecast,
		"Graft":                   StaticKeywordGraft,
		"Recover":                 StaticKeywordRecover,
		"Ripple":                  StaticKeywordRipple,
		"Split Second":            StaticKeywordSplitSecond,
		"Suspend":                 StaticKeywordSuspend,
		"Vanishing":               StaticKeywordVanishing,
		"Absorb":                  StaticKeywordAbsorb,
		"Aura Swap":               StaticKeywordAuraSwap,
		"Delve":                   StaticKeywordDelve,
		"Fortify":                 StaticKeywordFortify,
		"Frenzy":                  StaticKeywordFrenzy,
		"Gravestorm":              StaticKeywordGravestorm,
		"Poisonous":               StaticKeywordPoisonous,
		"Transfigure":             StaticKeywordTransfigure,
		"Champion":                StaticKeywordChampion,
		"Changeling":              StaticKeywordChangeling,
		"Evoke":                   StaticKeywordEvoke,
		"Hideaway":                StaticKeywordHideaway,
		"Prowl":                   StaticKeywordProwl,
		"Reinforce":               StaticKeywordReinforce,
		"Conspire":                StaticKeywordConspire,
		"Persist":                 StaticKeywordPersist,
		"Wither":                  StaticKeywordWither,
		"Retrace":                 StaticKeywordRetrace,
		"Devour":                  StaticKeywordDevour,
		"Exalted":                 StaticKeywordExalted,
		"Unearth":                 StaticKeywordUnearth,
		"Cascade":                 StaticKeywordCascade,
		"Annihilator":             StaticKeywordAnnihilator,
		"Level Up":                StaticKeywordLevelUp,
		"Rebound":                 StaticKeywordRebound,
		"Umbra Armor":             StaticKeywordUmbraArmor,
		"Infect":                  StaticKeywordInfect,
		"Battle Cry":              StaticKeywordBattleCry,
		"Living Weapon":           StaticKeywordLivingWeapon,
		"Undying":                 StaticKeywordUndying,
		"Miracle":                 StaticKeywordMiracle,
		"Soulbond":                StaticKeywordSoulbond,
		"Overload":                StaticKeywordOverload,
		"Scavenge":                StaticKeywordScavenge,
		"Unleash":                 StaticKeywordUnleash,
		"Cipher":                  StaticKeywordCipher,
		"Evolve":                  StaticKeywordEvolve,
		"Extort":                  StaticKeywordExtort,
		"Fuse":                    StaticKeywordFuse,
		"Bestow":                  StaticKeywordBestow,
		"Tribute":                 StaticKeywordTribute,
		"Dethrone":                StaticKeywordDethrone,
		"Hidden Agenda":           StaticKeywordHiddenAgenda,
		"Outlast":                 StaticKeywordOutlast,
		"Prowess":                 StaticKeywordProwess,
		"Dash":                    StaticKeywordDash,
		"Exploit":                 StaticKeywordExploit,
		"Menace":                  StaticKeywordMenace,
		"Renown":                  StaticKeywordRenown,
		"Awaken":                  StaticKeywordAwaken,
		"Devoid":                  StaticKeywordDevoid,
		"Ingest":                  StaticKeywordIngest,
		"Myriad":                  StaticKeywordMyriad,
		"Surge":                   StaticKeywordSurge,
		"Skulk":                   StaticKeywordSkulk,
		"Emerge":                  StaticKeywordEmerge,
		"Escalate":                StaticKeywordEscalate,
		"Melee":                   StaticKeywordMelee,
		"Crew":                    StaticKeywordCrew,
		"Fabricate":               StaticKeywordFabricate,
		"Partner":                 StaticKeywordPartner,
		"Undaunted":               StaticKeywordUndaunted,
		"Improvise":               StaticKeywordImprovise,
		"Aftermath":               StaticKeywordAftermath,
		"Embalm":                  StaticKeywordEmbalm,
		"Eternalize":              StaticKeywordEternalize,
		"Afflict":                 StaticKeywordAfflict,
		"Ascend":                  StaticKeywordAscend,
		"Assist":                  StaticKeywordAssist,
		"Jump-Start":              StaticKeywordJumpStart,
		"Mentor":                  StaticKeywordMentor,
		"Afterlife":               StaticKeywordAfterlife,
		"Riot":                    StaticKeywordRiot,
		"Spectacle":               StaticKeywordSpectacle,
		"Escape":                  StaticKeywordEscape,
		"Companion":               StaticKeywordCompanion,
		"Mutate":                  StaticKeywordMutate,
		"Encore":                  StaticKeywordEncore,
		"Boast":                   StaticKeywordBoast,
		"Foretell":                StaticKeywordForetell,
		"Demonstrate":             StaticKeywordDemonstrate,
		"Daybound and Nightbound": StaticKeywordDayboundAndNightbound,
		"Disturb":                 StaticKeywordDisturb,
		"Decayed":                 StaticKeywordDecayed,
		"Cleave":                  StaticKeywordCleave,
		"Training":                StaticKeywordTraining,
		"Compleated":              StaticKeywordCompleated,
		"Reconfigure":             StaticKeywordReconfigure,
		"Blitz":                   StaticKeywordBlitz,
		"Casualty":                StaticKeywordCasualty,
		"Enlist":                  StaticKeywordEnlist,
		"Read Ahead":              StaticKeywordReadAhead,
		"Ravenous":                StaticKeywordRavenous,
		"Squad":                   StaticKeywordSquad,
		"Space Sculptor":          StaticKeywordSpaceSculptor,
		"Visit":                   StaticKeywordVisit,
		"Prototype":               StaticKeywordPrototype,
		"Living Metal":            StaticKeywordLivingMetal,
		"More Than Meets the Eye": StaticKeywordMoreThanMeetsTheEye,
		"For Mirrodin!":           StaticKeywordForMirrodin,
		"Toxic":                   StaticKeywordToxic,
		"Backup":                  StaticKeywordBackup,
		"Bargain":                 StaticKeywordBargain,
		"Craft":                   StaticKeywordCraft,
		"Disguise":                StaticKeywordDisguise,
		"Solved":                  StaticKeywordSolved,
		"Plot":                    StaticKeywordPlot,
		"Saddle":                  StaticKeywordSaddle,
		"Spree":                   StaticKeywordSpree,
		"Freerunning":             StaticKeywordFreerunning,
		"Gift":                    StaticKeywordGift,
		"Offspring":               StaticKeywordOffspring,
		"Impending":               StaticKeywordImpending,
	}
	kw, ok := allStaticKeywords[s]
	return kw, ok
}

func (s StaticKeyword) Name() string {
	return string(s)
}
