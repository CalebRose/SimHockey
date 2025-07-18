/* Do not change, this code is generated from Golang structs */


export class BasePlayer {
    FirstName: string;
    LastName: string;
    Position: string;
    Archetype: string;
    TeamID: number;
    Team: string;
    Height: number;
    Weight: number;
    Stars: number;
    Age: number;
    Overall: number;
    Agility: number;
    Faceoffs: number;
    LongShotAccuracy: number;
    LongShotPower: number;
    CloseShotAccuracy: number;
    CloseShotPower: number;
    OneTimer: number;
    Passing: number;
    PuckHandling: number;
    Strength: number;
    BodyChecking: number;
    StickChecking: number;
    ShotBlocking: number;
    Goalkeeping: number;
    GoalieVision: number;
    GoalieReboundControl: number;
    Discipline: number;
    Aggression: number;
    Stamina: number;
    GoalieStamina: number;
    InjuryRating: number;
    DisciplineDeviation: number;
    InjuryDeviation: number;
    PrimeAge: number;
    Clutch: number;
    HighSchool: string;
    City: string;
    State: string;
    Country: string;
    OriginalTeamID: number;
    OriginalTeam: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Competitiveness: number;
    TeamLoyalty: number;
    PlaytimePreference: number;
    PlayerMorale: number;
    Personality: string;
    RelativeID: number;
    RelativeType: number;
    Notes: string;
    HasProgressed: boolean;
    ProgramPref: number;
    ProfDevPref: number;
    TraditionsPref: number;
    FacilitiesPref: number;
    AtmospherePref: number;
    AcademicsPref: number;
    ConferencePref: number;
    CoachPref: number;
    SeasonMomentumPref: number;
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.FirstName = source["FirstName"];
        this.LastName = source["LastName"];
        this.Position = source["Position"];
        this.Archetype = source["Archetype"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.Height = source["Height"];
        this.Weight = source["Weight"];
        this.Stars = source["Stars"];
        this.Age = source["Age"];
        this.Overall = source["Overall"];
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.GoalieReboundControl = source["GoalieReboundControl"];
        this.Discipline = source["Discipline"];
        this.Aggression = source["Aggression"];
        this.Stamina = source["Stamina"];
        this.GoalieStamina = source["GoalieStamina"];
        this.InjuryRating = source["InjuryRating"];
        this.DisciplineDeviation = source["DisciplineDeviation"];
        this.InjuryDeviation = source["InjuryDeviation"];
        this.PrimeAge = source["PrimeAge"];
        this.Clutch = source["Clutch"];
        this.HighSchool = source["HighSchool"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.OriginalTeamID = source["OriginalTeamID"];
        this.OriginalTeam = source["OriginalTeam"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Competitiveness = source["Competitiveness"];
        this.TeamLoyalty = source["TeamLoyalty"];
        this.PlaytimePreference = source["PlaytimePreference"];
        this.PlayerMorale = source["PlayerMorale"];
        this.Personality = source["Personality"];
        this.RelativeID = source["RelativeID"];
        this.RelativeType = source["RelativeType"];
        this.Notes = source["Notes"];
        this.HasProgressed = source["HasProgressed"];
        this.ProgramPref = source["ProgramPref"];
        this.ProfDevPref = source["ProfDevPref"];
        this.TraditionsPref = source["TraditionsPref"];
        this.FacilitiesPref = source["FacilitiesPref"];
        this.AtmospherePref = source["AtmospherePref"];
        this.AcademicsPref = source["AcademicsPref"];
        this.ConferencePref = source["ConferencePref"];
        this.CoachPref = source["CoachPref"];
        this.SeasonMomentumPref = source["SeasonMomentumPref"];
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
    }
}
export class DeletedAt {
    Time: Time;
    Valid: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Time = this.convertValues(source["Time"], Time);
        this.Valid = source["Valid"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Time {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class CollegeGameplan {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    IsAI: boolean;
    ForwardShotPreference: number;
    DefenderShotPreference: number;
    ForwardCheckPreference: number;
    DefenderCheckPreference: number;
    CenterSortPreference1: number;
    CenterSortPreference2: number;
    CenterSortPreference3: number;
    ForwardSortPreference1: number;
    ForwardSortPreference2: number;
    ForwardSortPreference3: number;
    DefenderSortPreference1: number;
    DefenderSortPreference2: number;
    DefenderSortPreference3: number;
    GoalieSortPreference: number;
    LongerPassesEnabled: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.IsAI = source["IsAI"];
        this.ForwardShotPreference = source["ForwardShotPreference"];
        this.DefenderShotPreference = source["DefenderShotPreference"];
        this.ForwardCheckPreference = source["ForwardCheckPreference"];
        this.DefenderCheckPreference = source["DefenderCheckPreference"];
        this.CenterSortPreference1 = source["CenterSortPreference1"];
        this.CenterSortPreference2 = source["CenterSortPreference2"];
        this.CenterSortPreference3 = source["CenterSortPreference3"];
        this.ForwardSortPreference1 = source["ForwardSortPreference1"];
        this.ForwardSortPreference2 = source["ForwardSortPreference2"];
        this.ForwardSortPreference3 = source["ForwardSortPreference3"];
        this.DefenderSortPreference1 = source["DefenderSortPreference1"];
        this.DefenderSortPreference2 = source["DefenderSortPreference2"];
        this.DefenderSortPreference3 = source["DefenderSortPreference3"];
        this.GoalieSortPreference = source["GoalieSortPreference"];
        this.LongerPassesEnabled = source["LongerPassesEnabled"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProGameplan {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    IsAI: boolean;
    ForwardShotPreference: number;
    DefenderShotPreference: number;
    ForwardCheckPreference: number;
    DefenderCheckPreference: number;
    CenterSortPreference1: number;
    CenterSortPreference2: number;
    CenterSortPreference3: number;
    ForwardSortPreference1: number;
    ForwardSortPreference2: number;
    ForwardSortPreference3: number;
    DefenderSortPreference1: number;
    DefenderSortPreference2: number;
    DefenderSortPreference3: number;
    GoalieSortPreference: number;
    LongerPassesEnabled: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.IsAI = source["IsAI"];
        this.ForwardShotPreference = source["ForwardShotPreference"];
        this.DefenderShotPreference = source["DefenderShotPreference"];
        this.ForwardCheckPreference = source["ForwardCheckPreference"];
        this.DefenderCheckPreference = source["DefenderCheckPreference"];
        this.CenterSortPreference1 = source["CenterSortPreference1"];
        this.CenterSortPreference2 = source["CenterSortPreference2"];
        this.CenterSortPreference3 = source["CenterSortPreference3"];
        this.ForwardSortPreference1 = source["ForwardSortPreference1"];
        this.ForwardSortPreference2 = source["ForwardSortPreference2"];
        this.ForwardSortPreference3 = source["ForwardSortPreference3"];
        this.DefenderSortPreference1 = source["DefenderSortPreference1"];
        this.DefenderSortPreference2 = source["DefenderSortPreference2"];
        this.DefenderSortPreference3 = source["DefenderSortPreference3"];
        this.GoalieSortPreference = source["GoalieSortPreference"];
        this.LongerPassesEnabled = source["LongerPassesEnabled"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class DraftPick {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    Season: number;
    DrafteeID: number;
    DraftRound: number;
    DraftNumber: number;
    TeamID: number;
    Team: string;
    OriginalTeamID: number;
    OriginalTeam: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    DraftValue: number;
    Notes: string;
    SelectedPlayerID: number;
    SelectedPlayerName: string;
    SelectedPlayerPosition: string;
    IsCompensation: boolean;
    IsVoid: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.Season = source["Season"];
        this.DrafteeID = source["DrafteeID"];
        this.DraftRound = source["DraftRound"];
        this.DraftNumber = source["DraftNumber"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.OriginalTeamID = source["OriginalTeamID"];
        this.OriginalTeam = source["OriginalTeam"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.DraftValue = source["DraftValue"];
        this.Notes = source["Notes"];
        this.SelectedPlayerID = source["SelectedPlayerID"];
        this.SelectedPlayerName = source["SelectedPlayerName"];
        this.SelectedPlayerPosition = source["SelectedPlayerPosition"];
        this.IsCompensation = source["IsCompensation"];
        this.IsVoid = source["IsVoid"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class TradePreferences {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    Centers: boolean;
    CentersType: number;
    Forwards: boolean;
    ForwardsType: number;
    Defenders: boolean;
    DefendersType: number;
    Goalies: boolean;
    GoaliesType: number;
    DraftPicks: boolean;
    DraftPickType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.Centers = source["Centers"];
        this.CentersType = source["CentersType"];
        this.Forwards = source["Forwards"];
        this.ForwardsType = source["ForwardsType"];
        this.Defenders = source["Defenders"];
        this.DefendersType = source["DefendersType"];
        this.Goalies = source["Goalies"];
        this.GoaliesType = source["GoaliesType"];
        this.DraftPicks = source["DraftPicks"];
        this.DraftPickType = source["DraftPickType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class TradeOption {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TradeProposalID: number;
    TeamID: number;
    PlayerID: number;
    DraftPickID: number;
    OptionType: string;
    SalaryPercentage: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TradeProposalID = source["TradeProposalID"];
        this.TeamID = source["TeamID"];
        this.PlayerID = source["PlayerID"];
        this.DraftPickID = source["DraftPickID"];
        this.OptionType = source["OptionType"];
        this.SalaryPercentage = source["SalaryPercentage"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class TradeProposal {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    RecepientTeamID: number;
    IsTradeAccepted: boolean;
    IsTradeRejected: boolean;
    IsSynced: boolean;
    TeamTradeOptions: TradeOption[];
    RecepientTeamTradeOptions: TradeOption[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.RecepientTeamID = source["RecepientTeamID"];
        this.IsTradeAccepted = source["IsTradeAccepted"];
        this.IsTradeRejected = source["IsTradeRejected"];
        this.IsSynced = source["IsSynced"];
        this.TeamTradeOptions = this.convertValues(source["TeamTradeOptions"], TradeOption);
        this.RecepientTeamTradeOptions = this.convertValues(source["RecepientTeamTradeOptions"], TradeOption);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}

export class ScoreBoard {
    P1Home: number;
    P2Home: number;
    P3Home: number;
    OTHome: number;
    P1Away: number;
    P2Away: number;
    P3Away: number;
    OTAway: number;
    HomeShootoutScore: number;
    AwayShootoutScore: number;
    HomeOffensiveScheme: string;
    HomeDefensiveScheme: string;
    AwayOffensiveScheme: string;
    AwayDefensiveScheme: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.P1Home = source["P1Home"];
        this.P2Home = source["P2Home"];
        this.P3Home = source["P3Home"];
        this.OTHome = source["OTHome"];
        this.P1Away = source["P1Away"];
        this.P2Away = source["P2Away"];
        this.P3Away = source["P3Away"];
        this.OTAway = source["OTAway"];
        this.HomeShootoutScore = source["HomeShootoutScore"];
        this.AwayShootoutScore = source["AwayShootoutScore"];
        this.HomeOffensiveScheme = source["HomeOffensiveScheme"];
        this.HomeDefensiveScheme = source["HomeDefensiveScheme"];
        this.AwayOffensiveScheme = source["AwayOffensiveScheme"];
        this.AwayDefensiveScheme = source["AwayDefensiveScheme"];
    }
}
export class ProfessionalPlayerGameStats {
    WeekID: number;
    GameID: number;
    RevealResults: boolean;
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    StartedGame: boolean;
    GameDay: string;
    PlayerID: number;
    TeamID: number;
    SeasonID: number;
    Team: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Goals: number;
    Assists: number;
    Points: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    TimeOnIce: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    GoalieWins: number;
    GoalieLosses: number;
    GoalieTies: number;
    OvertimeLosses: number;
    ShotsAgainst: number;
    Saves: number;
    GoalsAgainst: number;
    SavePercentage: number;
    Shutouts: number;
    ShotsBlocked: number;
    BodyChecks: number;
    StickChecks: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    GameType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.WeekID = source["WeekID"];
        this.GameID = source["GameID"];
        this.RevealResults = source["RevealResults"];
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.StartedGame = source["StartedGame"];
        this.GameDay = source["GameDay"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.SeasonID = source["SeasonID"];
        this.Team = source["Team"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Goals = source["Goals"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.TimeOnIce = source["TimeOnIce"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.GoalieWins = source["GoalieWins"];
        this.GoalieLosses = source["GoalieLosses"];
        this.GoalieTies = source["GoalieTies"];
        this.OvertimeLosses = source["OvertimeLosses"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.ShotsBlocked = source["ShotsBlocked"];
        this.BodyChecks = source["BodyChecks"];
        this.StickChecks = source["StickChecks"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.GameType = source["GameType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class PlayByPlayResponse {
    GameID: number;
    PlayNumber: number;
    HomeTeamID: number;
    HomeTeamScore: number;
    AwayTeamID: number;
    AwayTeamScore: number;
    Period: number;
    TimeOnClock: string;
    SecondsConsumed: number;
    Event: string;
    Zone: string;
    NextZone: string;
    Outcome: string;
    TeamID: number;
    PuckCarrierID: number;
    PassedPlayerID: number;
    AssistingPlayerID: number;
    DefenderID: number;
    GoalieID: number;
    InjuryID: number;
    InjuryType: number;
    InjuryDuration: number;
    Penalty: string;
    Severity: string;
    IsFight: string;
    IsBreakaway: boolean;
    Result: string;
    StreamResult: string[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.GameID = source["GameID"];
        this.PlayNumber = source["PlayNumber"];
        this.HomeTeamID = source["HomeTeamID"];
        this.HomeTeamScore = source["HomeTeamScore"];
        this.AwayTeamID = source["AwayTeamID"];
        this.AwayTeamScore = source["AwayTeamScore"];
        this.Period = source["Period"];
        this.TimeOnClock = source["TimeOnClock"];
        this.SecondsConsumed = source["SecondsConsumed"];
        this.Event = source["Event"];
        this.Zone = source["Zone"];
        this.NextZone = source["NextZone"];
        this.Outcome = source["Outcome"];
        this.TeamID = source["TeamID"];
        this.PuckCarrierID = source["PuckCarrierID"];
        this.PassedPlayerID = source["PassedPlayerID"];
        this.AssistingPlayerID = source["AssistingPlayerID"];
        this.DefenderID = source["DefenderID"];
        this.GoalieID = source["GoalieID"];
        this.InjuryID = source["InjuryID"];
        this.InjuryType = source["InjuryType"];
        this.InjuryDuration = source["InjuryDuration"];
        this.Penalty = source["Penalty"];
        this.Severity = source["Severity"];
        this.IsFight = source["IsFight"];
        this.IsBreakaway = source["IsBreakaway"];
        this.Result = source["Result"];
        this.StreamResult = source["StreamResult"];
    }
}
export class CollegePlayerGameStats {
    WeekID: number;
    GameID: number;
    RevealResults: boolean;
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    StartedGame: boolean;
    GameDay: string;
    PlayerID: number;
    TeamID: number;
    SeasonID: number;
    Team: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Goals: number;
    Assists: number;
    Points: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    TimeOnIce: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    GoalieWins: number;
    GoalieLosses: number;
    GoalieTies: number;
    OvertimeLosses: number;
    ShotsAgainst: number;
    Saves: number;
    GoalsAgainst: number;
    SavePercentage: number;
    Shutouts: number;
    ShotsBlocked: number;
    BodyChecks: number;
    StickChecks: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    GameType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.WeekID = source["WeekID"];
        this.GameID = source["GameID"];
        this.RevealResults = source["RevealResults"];
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.StartedGame = source["StartedGame"];
        this.GameDay = source["GameDay"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.SeasonID = source["SeasonID"];
        this.Team = source["Team"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Goals = source["Goals"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.TimeOnIce = source["TimeOnIce"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.GoalieWins = source["GoalieWins"];
        this.GoalieLosses = source["GoalieLosses"];
        this.GoalieTies = source["GoalieTies"];
        this.OvertimeLosses = source["OvertimeLosses"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.ShotsBlocked = source["ShotsBlocked"];
        this.BodyChecks = source["BodyChecks"];
        this.StickChecks = source["StickChecks"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.GameType = source["GameType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class GameResultsResponse {
    CHLHomeStats: CollegePlayerGameStats[];
    CHLAwayStats: CollegePlayerGameStats[];
    CHLPlayByPlays: PlayByPlayResponse[];
    PHLHomeStats: ProfessionalPlayerGameStats[];
    PHLAwayStats: ProfessionalPlayerGameStats[];
    PHLPlayByPlays: PlayByPlayResponse[];
    Score: ScoreBoard;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.CHLHomeStats = this.convertValues(source["CHLHomeStats"], CollegePlayerGameStats);
        this.CHLAwayStats = this.convertValues(source["CHLAwayStats"], CollegePlayerGameStats);
        this.CHLPlayByPlays = this.convertValues(source["CHLPlayByPlays"], PlayByPlayResponse);
        this.PHLHomeStats = this.convertValues(source["PHLHomeStats"], ProfessionalPlayerGameStats);
        this.PHLAwayStats = this.convertValues(source["PHLAwayStats"], ProfessionalPlayerGameStats);
        this.PHLPlayByPlays = this.convertValues(source["PHLPlayByPlays"], PlayByPlayResponse);
        this.Score = this.convertValues(source["Score"], ScoreBoard);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalTeamSeasonStats {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    TeamID: number;
    Team: string;
    GameDay: string;
    GoalsFor: number;
    GoalsAgainst: number;
    Assists: number;
    Points: number;
    Period1Score: number;
    Period2Score: number;
    Period3Score: number;
    OTScore: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    ShotsAgainst: number;
    Saves: number;
    SavePercentage: number;
    Shutouts: number;
    GameType: number;
    StatType: number;
    GamesPlayed: number;
    PointPercentage: number;
    RegulationWins: number;
    RegulationOvertimeWins: number;
    ShootoutsWon: number;
    GFGP: number;
    GAGP: number;
    PowerPlayPercentage: number;
    PenaltyKillPercentage: number;
    PowerPlayNetPercentage: number;
    PenaltyKillNetPercentage: number;
    ShotsPerGame: number;
    ShotsAgainstPerGame: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.GameDay = source["GameDay"];
        this.GoalsFor = source["GoalsFor"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.Period1Score = source["Period1Score"];
        this.Period2Score = source["Period2Score"];
        this.Period3Score = source["Period3Score"];
        this.OTScore = source["OTScore"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.GameType = source["GameType"];
        this.StatType = source["StatType"];
        this.GamesPlayed = source["GamesPlayed"];
        this.PointPercentage = source["PointPercentage"];
        this.RegulationWins = source["RegulationWins"];
        this.RegulationOvertimeWins = source["RegulationOvertimeWins"];
        this.ShootoutsWon = source["ShootoutsWon"];
        this.GFGP = source["GFGP"];
        this.GAGP = source["GAGP"];
        this.PowerPlayPercentage = source["PowerPlayPercentage"];
        this.PenaltyKillPercentage = source["PenaltyKillPercentage"];
        this.PowerPlayNetPercentage = source["PowerPlayNetPercentage"];
        this.PenaltyKillNetPercentage = source["PenaltyKillNetPercentage"];
        this.ShotsPerGame = source["ShotsPerGame"];
        this.ShotsAgainstPerGame = source["ShotsAgainstPerGame"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalTeamGameStats {
    WeekID: number;
    GameID: number;
    RevealResults: boolean;
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    TeamID: number;
    Team: string;
    GameDay: string;
    GoalsFor: number;
    GoalsAgainst: number;
    Assists: number;
    Points: number;
    Period1Score: number;
    Period2Score: number;
    Period3Score: number;
    OTScore: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    ShotsAgainst: number;
    Saves: number;
    SavePercentage: number;
    Shutouts: number;
    GameType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.WeekID = source["WeekID"];
        this.GameID = source["GameID"];
        this.RevealResults = source["RevealResults"];
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.GameDay = source["GameDay"];
        this.GoalsFor = source["GoalsFor"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.Period1Score = source["Period1Score"];
        this.Period2Score = source["Period2Score"];
        this.Period3Score = source["Period3Score"];
        this.OTScore = source["OTScore"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.GameType = source["GameType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalPlayerSeasonStats {
    StatType: number;
    GamesPlayed: number;
    GamesStarted: number;
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    StartedGame: boolean;
    GameDay: string;
    PlayerID: number;
    TeamID: number;
    SeasonID: number;
    Team: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Goals: number;
    Assists: number;
    Points: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    TimeOnIce: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    GoalieWins: number;
    GoalieLosses: number;
    GoalieTies: number;
    OvertimeLosses: number;
    ShotsAgainst: number;
    Saves: number;
    GoalsAgainst: number;
    SavePercentage: number;
    Shutouts: number;
    ShotsBlocked: number;
    BodyChecks: number;
    StickChecks: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    GameType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.StatType = source["StatType"];
        this.GamesPlayed = source["GamesPlayed"];
        this.GamesStarted = source["GamesStarted"];
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.StartedGame = source["StartedGame"];
        this.GameDay = source["GameDay"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.SeasonID = source["SeasonID"];
        this.Team = source["Team"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Goals = source["Goals"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.TimeOnIce = source["TimeOnIce"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.GoalieWins = source["GoalieWins"];
        this.GoalieLosses = source["GoalieLosses"];
        this.GoalieTies = source["GoalieTies"];
        this.OvertimeLosses = source["OvertimeLosses"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.ShotsBlocked = source["ShotsBlocked"];
        this.BodyChecks = source["BodyChecks"];
        this.StickChecks = source["StickChecks"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.GameType = source["GameType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegeTeamSeasonStats {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    TeamID: number;
    Team: string;
    GameDay: string;
    GoalsFor: number;
    GoalsAgainst: number;
    Assists: number;
    Points: number;
    Period1Score: number;
    Period2Score: number;
    Period3Score: number;
    OTScore: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    ShotsAgainst: number;
    Saves: number;
    SavePercentage: number;
    Shutouts: number;
    GameType: number;
    StatType: number;
    GamesPlayed: number;
    PointPercentage: number;
    RegulationWins: number;
    RegulationOvertimeWins: number;
    ShootoutsWon: number;
    GFGP: number;
    GAGP: number;
    PowerPlayPercentage: number;
    PenaltyKillPercentage: number;
    PowerPlayNetPercentage: number;
    PenaltyKillNetPercentage: number;
    ShotsPerGame: number;
    ShotsAgainstPerGame: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.GameDay = source["GameDay"];
        this.GoalsFor = source["GoalsFor"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.Period1Score = source["Period1Score"];
        this.Period2Score = source["Period2Score"];
        this.Period3Score = source["Period3Score"];
        this.OTScore = source["OTScore"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.GameType = source["GameType"];
        this.StatType = source["StatType"];
        this.GamesPlayed = source["GamesPlayed"];
        this.PointPercentage = source["PointPercentage"];
        this.RegulationWins = source["RegulationWins"];
        this.RegulationOvertimeWins = source["RegulationOvertimeWins"];
        this.ShootoutsWon = source["ShootoutsWon"];
        this.GFGP = source["GFGP"];
        this.GAGP = source["GAGP"];
        this.PowerPlayPercentage = source["PowerPlayPercentage"];
        this.PenaltyKillPercentage = source["PenaltyKillPercentage"];
        this.PowerPlayNetPercentage = source["PowerPlayNetPercentage"];
        this.PenaltyKillNetPercentage = source["PenaltyKillNetPercentage"];
        this.ShotsPerGame = source["ShotsPerGame"];
        this.ShotsAgainstPerGame = source["ShotsAgainstPerGame"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegeTeamGameStats {
    WeekID: number;
    GameID: number;
    RevealResults: boolean;
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    TeamID: number;
    Team: string;
    GameDay: string;
    GoalsFor: number;
    GoalsAgainst: number;
    Assists: number;
    Points: number;
    Period1Score: number;
    Period2Score: number;
    Period3Score: number;
    OTScore: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    ShotsAgainst: number;
    Saves: number;
    SavePercentage: number;
    Shutouts: number;
    GameType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.WeekID = source["WeekID"];
        this.GameID = source["GameID"];
        this.RevealResults = source["RevealResults"];
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.GameDay = source["GameDay"];
        this.GoalsFor = source["GoalsFor"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.Period1Score = source["Period1Score"];
        this.Period2Score = source["Period2Score"];
        this.Period3Score = source["Period3Score"];
        this.OTScore = source["OTScore"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.GameType = source["GameType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegePlayerSeasonStats {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    StartedGame: boolean;
    GameDay: string;
    PlayerID: number;
    TeamID: number;
    SeasonID: number;
    Team: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Goals: number;
    Assists: number;
    Points: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    TimeOnIce: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    GoalieWins: number;
    GoalieLosses: number;
    GoalieTies: number;
    OvertimeLosses: number;
    ShotsAgainst: number;
    Saves: number;
    GoalsAgainst: number;
    SavePercentage: number;
    Shutouts: number;
    ShotsBlocked: number;
    BodyChecks: number;
    StickChecks: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    GameType: number;
    StatType: number;
    GamesPlayed: number;
    GamesStarted: number;
    PointsPerGamePlayed: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.StartedGame = source["StartedGame"];
        this.GameDay = source["GameDay"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.SeasonID = source["SeasonID"];
        this.Team = source["Team"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Goals = source["Goals"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.TimeOnIce = source["TimeOnIce"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.GoalieWins = source["GoalieWins"];
        this.GoalieLosses = source["GoalieLosses"];
        this.GoalieTies = source["GoalieTies"];
        this.OvertimeLosses = source["OvertimeLosses"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.ShotsBlocked = source["ShotsBlocked"];
        this.BodyChecks = source["BodyChecks"];
        this.StickChecks = source["StickChecks"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.GameType = source["GameType"];
        this.StatType = source["StatType"];
        this.GamesPlayed = source["GamesPlayed"];
        this.GamesStarted = source["GamesStarted"];
        this.PointsPerGamePlayed = source["PointsPerGamePlayed"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class SearchStatsResponse {
    CHLPlayerGameStats: CollegePlayerGameStats[];
    CHLPlayerSeasonStats: CollegePlayerSeasonStats[];
    CHLTeamGameStats: CollegeTeamGameStats[];
    CHLTeamSeasonStats: CollegeTeamSeasonStats[];
    PHLPlayerGameStats: ProfessionalPlayerGameStats[];
    PHLPlayerSeasonStats: ProfessionalPlayerSeasonStats[];
    PHLTeamGameStats: ProfessionalTeamGameStats[];
    PHLTeamSeasonStats: ProfessionalTeamSeasonStats[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.CHLPlayerGameStats = this.convertValues(source["CHLPlayerGameStats"], CollegePlayerGameStats);
        this.CHLPlayerSeasonStats = this.convertValues(source["CHLPlayerSeasonStats"], CollegePlayerSeasonStats);
        this.CHLTeamGameStats = this.convertValues(source["CHLTeamGameStats"], CollegeTeamGameStats);
        this.CHLTeamSeasonStats = this.convertValues(source["CHLTeamSeasonStats"], CollegeTeamSeasonStats);
        this.PHLPlayerGameStats = this.convertValues(source["PHLPlayerGameStats"], ProfessionalPlayerGameStats);
        this.PHLPlayerSeasonStats = this.convertValues(source["PHLPlayerSeasonStats"], ProfessionalPlayerSeasonStats);
        this.PHLTeamGameStats = this.convertValues(source["PHLTeamGameStats"], ProfessionalTeamGameStats);
        this.PHLTeamSeasonStats = this.convertValues(source["PHLTeamSeasonStats"], ProfessionalTeamSeasonStats);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class PlayerPreferences {
    ProgramPref: number;
    ProfDevPref: number;
    TraditionsPref: number;
    FacilitiesPref: number;
    AtmospherePref: number;
    AcademicsPref: number;
    ConferencePref: number;
    CoachPref: number;
    SeasonMomentumPref: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ProgramPref = source["ProgramPref"];
        this.ProfDevPref = source["ProfDevPref"];
        this.TraditionsPref = source["TraditionsPref"];
        this.FacilitiesPref = source["FacilitiesPref"];
        this.AtmospherePref = source["AtmospherePref"];
        this.AcademicsPref = source["AcademicsPref"];
        this.ConferencePref = source["ConferencePref"];
        this.CoachPref = source["CoachPref"];
        this.SeasonMomentumPref = source["SeasonMomentumPref"];
    }
}
export class BasePlayerProgressions {
    Agility: number;
    Faceoffs: number;
    LongShotAccuracy: number;
    LongShotPower: number;
    CloseShotAccuracy: number;
    CloseShotPower: number;
    OneTimer: number;
    Passing: number;
    PuckHandling: number;
    Strength: number;
    BodyChecking: number;
    StickChecking: number;
    ShotBlocking: number;
    Goalkeeping: number;
    GoalieVision: number;
    GoalieReboundControl: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.GoalieReboundControl = source["GoalieReboundControl"];
    }
}
export class BasePotentials {
    AgilityPotential: number;
    FaceoffsPotential: number;
    CloseShotAccuracyPotential: number;
    CloseShotPowerPotential: number;
    LongShotAccuracyPotential: number;
    LongShotPowerPotential: number;
    PassingPotential: number;
    PuckHandlingPotential: number;
    StrengthPotential: number;
    BodyCheckingPotential: number;
    StickCheckingPotential: number;
    ShotBlockingPotential: number;
    GoalkeepingPotential: number;
    GoalieVisionPotential: number;
    GoalieReboundPotential: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.AgilityPotential = source["AgilityPotential"];
        this.FaceoffsPotential = source["FaceoffsPotential"];
        this.CloseShotAccuracyPotential = source["CloseShotAccuracyPotential"];
        this.CloseShotPowerPotential = source["CloseShotPowerPotential"];
        this.LongShotAccuracyPotential = source["LongShotAccuracyPotential"];
        this.LongShotPowerPotential = source["LongShotPowerPotential"];
        this.PassingPotential = source["PassingPotential"];
        this.PuckHandlingPotential = source["PuckHandlingPotential"];
        this.StrengthPotential = source["StrengthPotential"];
        this.BodyCheckingPotential = source["BodyCheckingPotential"];
        this.StickCheckingPotential = source["StickCheckingPotential"];
        this.ShotBlockingPotential = source["ShotBlockingPotential"];
        this.GoalkeepingPotential = source["GoalkeepingPotential"];
        this.GoalieVisionPotential = source["GoalieVisionPotential"];
        this.GoalieReboundPotential = source["GoalieReboundPotential"];
    }
}
export class BaseInjuryData {
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    InjuryCount: number;
    Regression: number;
    DecayRate: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.InjuryCount = source["InjuryCount"];
        this.Regression = source["Regression"];
        this.DecayRate = source["DecayRate"];
    }
}
export class CollegePromise {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    CollegePlayerID: number;
    PromiseType: string;
    PromiseWeight: string;
    Benchmark: number;
    BenchmarkStr: string;
    PromiseMade: boolean;
    IsFullfilled: boolean;
    IsActive: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.CollegePlayerID = source["CollegePlayerID"];
        this.PromiseType = source["PromiseType"];
        this.PromiseWeight = source["PromiseWeight"];
        this.Benchmark = source["Benchmark"];
        this.BenchmarkStr = source["BenchmarkStr"];
        this.PromiseMade = source["PromiseMade"];
        this.IsFullfilled = source["IsFullfilled"];
        this.IsActive = source["IsActive"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class NullInt64 {
    Int64: number;
    Valid: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Int64 = source["Int64"];
        this.Valid = source["Valid"];
    }
}
export class TransferPortalProfile {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    CollegePlayerID: number;
    ProfileID: number;
    PromiseID: NullInt64;
    TeamAbbreviation: string;
    TotalPoints: number;
    CurrentWeeksPoints: number;
    PreviouslySpentPoints: number;
    SpendingCount: number;
    RemovedFromBoard: boolean;
    RolledOnPromise: boolean;
    LockProfile: boolean;
    IsSigned: boolean;
    Recruiter: string;
    CollegePlayer: CollegePlayer;
    Promise: CollegePromise;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.CollegePlayerID = source["CollegePlayerID"];
        this.ProfileID = source["ProfileID"];
        this.PromiseID = this.convertValues(source["PromiseID"], NullInt64);
        this.TeamAbbreviation = source["TeamAbbreviation"];
        this.TotalPoints = source["TotalPoints"];
        this.CurrentWeeksPoints = source["CurrentWeeksPoints"];
        this.PreviouslySpentPoints = source["PreviouslySpentPoints"];
        this.SpendingCount = source["SpendingCount"];
        this.RemovedFromBoard = source["RemovedFromBoard"];
        this.RolledOnPromise = source["RolledOnPromise"];
        this.LockProfile = source["LockProfile"];
        this.IsSigned = source["IsSigned"];
        this.Recruiter = source["Recruiter"];
        this.CollegePlayer = this.convertValues(source["CollegePlayer"], CollegePlayer);
        this.Promise = this.convertValues(source["Promise"], CollegePromise);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegePlayer {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    FirstName: string;
    LastName: string;
    Position: string;
    Archetype: string;
    TeamID: number;
    Team: string;
    Height: number;
    Weight: number;
    Stars: number;
    Age: number;
    Overall: number;
    Agility: number;
    Faceoffs: number;
    LongShotAccuracy: number;
    LongShotPower: number;
    CloseShotAccuracy: number;
    CloseShotPower: number;
    OneTimer: number;
    Passing: number;
    PuckHandling: number;
    Strength: number;
    BodyChecking: number;
    StickChecking: number;
    ShotBlocking: number;
    Goalkeeping: number;
    GoalieVision: number;
    GoalieReboundControl: number;
    Discipline: number;
    Aggression: number;
    Stamina: number;
    GoalieStamina: number;
    InjuryRating: number;
    DisciplineDeviation: number;
    InjuryDeviation: number;
    PrimeAge: number;
    Clutch: number;
    HighSchool: string;
    City: string;
    State: string;
    Country: string;
    OriginalTeamID: number;
    OriginalTeam: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Competitiveness: number;
    TeamLoyalty: number;
    PlaytimePreference: number;
    PlayerMorale: number;
    Personality: string;
    RelativeID: number;
    RelativeType: number;
    Notes: string;
    HasProgressed: boolean;
    ProgramPref: number;
    ProfDevPref: number;
    TraditionsPref: number;
    FacilitiesPref: number;
    AtmospherePref: number;
    AcademicsPref: number;
    ConferencePref: number;
    CoachPref: number;
    SeasonMomentumPref: number;
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;
    AgilityPotential: number;
    FaceoffsPotential: number;
    CloseShotAccuracyPotential: number;
    CloseShotPowerPotential: number;
    LongShotAccuracyPotential: number;
    LongShotPowerPotential: number;
    PassingPotential: number;
    PuckHandlingPotential: number;
    StrengthPotential: number;
    BodyCheckingPotential: number;
    StickCheckingPotential: number;
    ShotBlockingPotential: number;
    GoalkeepingPotential: number;
    GoalieVisionPotential: number;
    GoalieReboundPotential: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    InjuryCount: number;
    Regression: number;
    DecayRate: number;
    Year: number;
    IsRedshirt: boolean;
    IsRedshirting: boolean;
    TransferStatus: number;
    TransferLikeliness: string;
    DraftedTeamID: number;
    DraftedTeam: string;
    DraftedRound: number;
    DraftPickID: number;
    DraftedPick: number;
    Stats: CollegePlayerGameStats[];
    SeasonStats: CollegePlayerSeasonStats;
    Profiles: TransferPortalProfile[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.FirstName = source["FirstName"];
        this.LastName = source["LastName"];
        this.Position = source["Position"];
        this.Archetype = source["Archetype"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.Height = source["Height"];
        this.Weight = source["Weight"];
        this.Stars = source["Stars"];
        this.Age = source["Age"];
        this.Overall = source["Overall"];
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.GoalieReboundControl = source["GoalieReboundControl"];
        this.Discipline = source["Discipline"];
        this.Aggression = source["Aggression"];
        this.Stamina = source["Stamina"];
        this.GoalieStamina = source["GoalieStamina"];
        this.InjuryRating = source["InjuryRating"];
        this.DisciplineDeviation = source["DisciplineDeviation"];
        this.InjuryDeviation = source["InjuryDeviation"];
        this.PrimeAge = source["PrimeAge"];
        this.Clutch = source["Clutch"];
        this.HighSchool = source["HighSchool"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.OriginalTeamID = source["OriginalTeamID"];
        this.OriginalTeam = source["OriginalTeam"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Competitiveness = source["Competitiveness"];
        this.TeamLoyalty = source["TeamLoyalty"];
        this.PlaytimePreference = source["PlaytimePreference"];
        this.PlayerMorale = source["PlayerMorale"];
        this.Personality = source["Personality"];
        this.RelativeID = source["RelativeID"];
        this.RelativeType = source["RelativeType"];
        this.Notes = source["Notes"];
        this.HasProgressed = source["HasProgressed"];
        this.ProgramPref = source["ProgramPref"];
        this.ProfDevPref = source["ProfDevPref"];
        this.TraditionsPref = source["TraditionsPref"];
        this.FacilitiesPref = source["FacilitiesPref"];
        this.AtmospherePref = source["AtmospherePref"];
        this.AcademicsPref = source["AcademicsPref"];
        this.ConferencePref = source["ConferencePref"];
        this.CoachPref = source["CoachPref"];
        this.SeasonMomentumPref = source["SeasonMomentumPref"];
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
        this.AgilityPotential = source["AgilityPotential"];
        this.FaceoffsPotential = source["FaceoffsPotential"];
        this.CloseShotAccuracyPotential = source["CloseShotAccuracyPotential"];
        this.CloseShotPowerPotential = source["CloseShotPowerPotential"];
        this.LongShotAccuracyPotential = source["LongShotAccuracyPotential"];
        this.LongShotPowerPotential = source["LongShotPowerPotential"];
        this.PassingPotential = source["PassingPotential"];
        this.PuckHandlingPotential = source["PuckHandlingPotential"];
        this.StrengthPotential = source["StrengthPotential"];
        this.BodyCheckingPotential = source["BodyCheckingPotential"];
        this.StickCheckingPotential = source["StickCheckingPotential"];
        this.ShotBlockingPotential = source["ShotBlockingPotential"];
        this.GoalkeepingPotential = source["GoalkeepingPotential"];
        this.GoalieVisionPotential = source["GoalieVisionPotential"];
        this.GoalieReboundPotential = source["GoalieReboundPotential"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.InjuryCount = source["InjuryCount"];
        this.Regression = source["Regression"];
        this.DecayRate = source["DecayRate"];
        this.Year = source["Year"];
        this.IsRedshirt = source["IsRedshirt"];
        this.IsRedshirting = source["IsRedshirting"];
        this.TransferStatus = source["TransferStatus"];
        this.TransferLikeliness = source["TransferLikeliness"];
        this.DraftedTeamID = source["DraftedTeamID"];
        this.DraftedTeam = source["DraftedTeam"];
        this.DraftedRound = source["DraftedRound"];
        this.DraftPickID = source["DraftPickID"];
        this.DraftedPick = source["DraftedPick"];
        this.Stats = this.convertValues(source["Stats"], CollegePlayerGameStats);
        this.SeasonStats = this.convertValues(source["SeasonStats"], CollegePlayerSeasonStats);
        this.Profiles = this.convertValues(source["Profiles"], TransferPortalProfile);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ExtensionOffer {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    PlayerID: number;
    TeamID: number;
    SeasonID: number;
    ContractLength: number;
    Y1BaseSalary: number;
    Y2BaseSalary: number;
    Y3BaseSalary: number;
    Y4BaseSalary: number;
    Y5BaseSalary: number;
    TotalSalary: number;
    ContractValue: number;
    BonusPercentage: number;
    Rejections: number;
    IsAccepted: boolean;
    IsActive: boolean;
    IsRejected: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.SeasonID = source["SeasonID"];
        this.ContractLength = source["ContractLength"];
        this.Y1BaseSalary = source["Y1BaseSalary"];
        this.Y2BaseSalary = source["Y2BaseSalary"];
        this.Y3BaseSalary = source["Y3BaseSalary"];
        this.Y4BaseSalary = source["Y4BaseSalary"];
        this.Y5BaseSalary = source["Y5BaseSalary"];
        this.TotalSalary = source["TotalSalary"];
        this.ContractValue = source["ContractValue"];
        this.BonusPercentage = source["BonusPercentage"];
        this.Rejections = source["Rejections"];
        this.IsAccepted = source["IsAccepted"];
        this.IsActive = source["IsActive"];
        this.IsRejected = source["IsRejected"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class WaiverOffer {
    ID: number;
    PlayerID: number;
    TeamID: number;
    Team: string;
    WaiverOrder: number;
    IsActive: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.WaiverOrder = source["WaiverOrder"];
        this.IsActive = source["IsActive"];
    }
}
export class FreeAgencyOffer {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    PlayerID: number;
    TeamID: number;
    ContractLength: number;
    Y1BaseSalary: number;
    Y2BaseSalary: number;
    Y3BaseSalary: number;
    Y4BaseSalary: number;
    Y5BaseSalary: number;
    TotalBonus: number;
    TotalSalary: number;
    ContractValue: number;
    BonusPercentage: number;
    IsActive: boolean;
    IsRejected: boolean;
    Syncs: number;
    ToAffiliate: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.ContractLength = source["ContractLength"];
        this.Y1BaseSalary = source["Y1BaseSalary"];
        this.Y2BaseSalary = source["Y2BaseSalary"];
        this.Y3BaseSalary = source["Y3BaseSalary"];
        this.Y4BaseSalary = source["Y4BaseSalary"];
        this.Y5BaseSalary = source["Y5BaseSalary"];
        this.TotalBonus = source["TotalBonus"];
        this.TotalSalary = source["TotalSalary"];
        this.ContractValue = source["ContractValue"];
        this.BonusPercentage = source["BonusPercentage"];
        this.IsActive = source["IsActive"];
        this.IsRejected = source["IsRejected"];
        this.Syncs = source["Syncs"];
        this.ToAffiliate = source["ToAffiliate"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProContract {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    PlayerID: number;
    TeamID: number;
    OriginalTeamID: number;
    ContractLength: number;
    Y1BaseSalary: number;
    Y2BaseSalary: number;
    Y3BaseSalary: number;
    Y4BaseSalary: number;
    Y5BaseSalary: number;
    BonusPercentage: number;
    ContractType: string;
    ContractValue: number;
    SigningValue: number;
    IsActive: boolean;
    IsComplete: boolean;
    IsExtended: boolean;
    HasProgressed: boolean;
    PlayerRetired: boolean;
    TagType: number;
    IsTagged: boolean;
    IsCut: boolean;
    NoTradeClause: boolean;
    NoMovementClause: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.OriginalTeamID = source["OriginalTeamID"];
        this.ContractLength = source["ContractLength"];
        this.Y1BaseSalary = source["Y1BaseSalary"];
        this.Y2BaseSalary = source["Y2BaseSalary"];
        this.Y3BaseSalary = source["Y3BaseSalary"];
        this.Y4BaseSalary = source["Y4BaseSalary"];
        this.Y5BaseSalary = source["Y5BaseSalary"];
        this.BonusPercentage = source["BonusPercentage"];
        this.ContractType = source["ContractType"];
        this.ContractValue = source["ContractValue"];
        this.SigningValue = source["SigningValue"];
        this.IsActive = source["IsActive"];
        this.IsComplete = source["IsComplete"];
        this.IsExtended = source["IsExtended"];
        this.HasProgressed = source["HasProgressed"];
        this.PlayerRetired = source["PlayerRetired"];
        this.TagType = source["TagType"];
        this.IsTagged = source["IsTagged"];
        this.IsCut = source["IsCut"];
        this.NoTradeClause = source["NoTradeClause"];
        this.NoMovementClause = source["NoMovementClause"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalPlayer {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    FirstName: string;
    LastName: string;
    Position: string;
    Archetype: string;
    TeamID: number;
    Team: string;
    Height: number;
    Weight: number;
    Stars: number;
    Age: number;
    Overall: number;
    Agility: number;
    Faceoffs: number;
    LongShotAccuracy: number;
    LongShotPower: number;
    CloseShotAccuracy: number;
    CloseShotPower: number;
    OneTimer: number;
    Passing: number;
    PuckHandling: number;
    Strength: number;
    BodyChecking: number;
    StickChecking: number;
    ShotBlocking: number;
    Goalkeeping: number;
    GoalieVision: number;
    GoalieReboundControl: number;
    Discipline: number;
    Aggression: number;
    Stamina: number;
    GoalieStamina: number;
    InjuryRating: number;
    DisciplineDeviation: number;
    InjuryDeviation: number;
    PrimeAge: number;
    Clutch: number;
    HighSchool: string;
    City: string;
    State: string;
    Country: string;
    OriginalTeamID: number;
    OriginalTeam: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Competitiveness: number;
    TeamLoyalty: number;
    PlaytimePreference: number;
    PlayerMorale: number;
    Personality: string;
    RelativeID: number;
    RelativeType: number;
    Notes: string;
    HasProgressed: boolean;
    ProgramPref: number;
    ProfDevPref: number;
    TraditionsPref: number;
    FacilitiesPref: number;
    AtmospherePref: number;
    AcademicsPref: number;
    ConferencePref: number;
    CoachPref: number;
    SeasonMomentumPref: number;
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;
    AgilityPotential: number;
    FaceoffsPotential: number;
    CloseShotAccuracyPotential: number;
    CloseShotPowerPotential: number;
    LongShotAccuracyPotential: number;
    LongShotPowerPotential: number;
    PassingPotential: number;
    PuckHandlingPotential: number;
    StrengthPotential: number;
    BodyCheckingPotential: number;
    StickCheckingPotential: number;
    ShotBlockingPotential: number;
    GoalkeepingPotential: number;
    GoalieVisionPotential: number;
    GoalieReboundPotential: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    InjuryCount: number;
    Regression: number;
    DecayRate: number;
    CollegeID: number;
    Year: number;
    IsAffiliatePlayer: boolean;
    IsWaived: boolean;
    IsFreeAgent: boolean;
    IsOnTradeBlock: boolean;
    IsAcceptingOffers: boolean;
    IsNegotiating: boolean;
    DraftedTeamID: number;
    DraftedTeam: string;
    DraftedRound: number;
    DraftPickID: number;
    DraftedPick: number;
    MinimumValue: number;
    HasProgressed: boolean;
    Rejections: number;
    AffiliateTeamID: number;
    Marketability: number;
    JerseyPrice: number;
    MarketPreference: number;
    CompetitivePreference: number;
    FinancialPreference: number;
    IsEligibleForPlay: boolean;
    Stats: ProfessionalPlayerGameStats[];
    SeasonStats: ProfessionalPlayerSeasonStats;
    Contract: ProContract;
    Offers: FreeAgencyOffer[];
    WaiverOffer: WaiverOffer[];
    Extensions: ExtensionOffer[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.FirstName = source["FirstName"];
        this.LastName = source["LastName"];
        this.Position = source["Position"];
        this.Archetype = source["Archetype"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.Height = source["Height"];
        this.Weight = source["Weight"];
        this.Stars = source["Stars"];
        this.Age = source["Age"];
        this.Overall = source["Overall"];
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.GoalieReboundControl = source["GoalieReboundControl"];
        this.Discipline = source["Discipline"];
        this.Aggression = source["Aggression"];
        this.Stamina = source["Stamina"];
        this.GoalieStamina = source["GoalieStamina"];
        this.InjuryRating = source["InjuryRating"];
        this.DisciplineDeviation = source["DisciplineDeviation"];
        this.InjuryDeviation = source["InjuryDeviation"];
        this.PrimeAge = source["PrimeAge"];
        this.Clutch = source["Clutch"];
        this.HighSchool = source["HighSchool"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.OriginalTeamID = source["OriginalTeamID"];
        this.OriginalTeam = source["OriginalTeam"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Competitiveness = source["Competitiveness"];
        this.TeamLoyalty = source["TeamLoyalty"];
        this.PlaytimePreference = source["PlaytimePreference"];
        this.PlayerMorale = source["PlayerMorale"];
        this.Personality = source["Personality"];
        this.RelativeID = source["RelativeID"];
        this.RelativeType = source["RelativeType"];
        this.Notes = source["Notes"];
        this.HasProgressed = source["HasProgressed"];
        this.ProgramPref = source["ProgramPref"];
        this.ProfDevPref = source["ProfDevPref"];
        this.TraditionsPref = source["TraditionsPref"];
        this.FacilitiesPref = source["FacilitiesPref"];
        this.AtmospherePref = source["AtmospherePref"];
        this.AcademicsPref = source["AcademicsPref"];
        this.ConferencePref = source["ConferencePref"];
        this.CoachPref = source["CoachPref"];
        this.SeasonMomentumPref = source["SeasonMomentumPref"];
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
        this.AgilityPotential = source["AgilityPotential"];
        this.FaceoffsPotential = source["FaceoffsPotential"];
        this.CloseShotAccuracyPotential = source["CloseShotAccuracyPotential"];
        this.CloseShotPowerPotential = source["CloseShotPowerPotential"];
        this.LongShotAccuracyPotential = source["LongShotAccuracyPotential"];
        this.LongShotPowerPotential = source["LongShotPowerPotential"];
        this.PassingPotential = source["PassingPotential"];
        this.PuckHandlingPotential = source["PuckHandlingPotential"];
        this.StrengthPotential = source["StrengthPotential"];
        this.BodyCheckingPotential = source["BodyCheckingPotential"];
        this.StickCheckingPotential = source["StickCheckingPotential"];
        this.ShotBlockingPotential = source["ShotBlockingPotential"];
        this.GoalkeepingPotential = source["GoalkeepingPotential"];
        this.GoalieVisionPotential = source["GoalieVisionPotential"];
        this.GoalieReboundPotential = source["GoalieReboundPotential"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.InjuryCount = source["InjuryCount"];
        this.Regression = source["Regression"];
        this.DecayRate = source["DecayRate"];
        this.CollegeID = source["CollegeID"];
        this.Year = source["Year"];
        this.IsAffiliatePlayer = source["IsAffiliatePlayer"];
        this.IsWaived = source["IsWaived"];
        this.IsFreeAgent = source["IsFreeAgent"];
        this.IsOnTradeBlock = source["IsOnTradeBlock"];
        this.IsAcceptingOffers = source["IsAcceptingOffers"];
        this.IsNegotiating = source["IsNegotiating"];
        this.DraftedTeamID = source["DraftedTeamID"];
        this.DraftedTeam = source["DraftedTeam"];
        this.DraftedRound = source["DraftedRound"];
        this.DraftPickID = source["DraftPickID"];
        this.DraftedPick = source["DraftedPick"];
        this.MinimumValue = source["MinimumValue"];
        this.HasProgressed = source["HasProgressed"];
        this.Rejections = source["Rejections"];
        this.AffiliateTeamID = source["AffiliateTeamID"];
        this.Marketability = source["Marketability"];
        this.JerseyPrice = source["JerseyPrice"];
        this.MarketPreference = source["MarketPreference"];
        this.CompetitivePreference = source["CompetitivePreference"];
        this.FinancialPreference = source["FinancialPreference"];
        this.IsEligibleForPlay = source["IsEligibleForPlay"];
        this.Stats = this.convertValues(source["Stats"], ProfessionalPlayerGameStats);
        this.SeasonStats = this.convertValues(source["SeasonStats"], ProfessionalPlayerSeasonStats);
        this.Contract = this.convertValues(source["Contract"], ProContract);
        this.Offers = this.convertValues(source["Offers"], FreeAgencyOffer);
        this.WaiverOffer = this.convertValues(source["WaiverOffer"], WaiverOffer);
        this.Extensions = this.convertValues(source["Extensions"], ExtensionOffer);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class BaseTeam {
    TeamName: string;
    Mascot: string;
    Abbreviation: string;
    ConferenceID: number;
    Conference: string;
    Coach: string;
    City: string;
    State: string;
    Country: string;
    ArenaID: number;
    Arena: string;
    ArenaCapacity: number;
    RecordAttendance: number;
    FirstPlayed: number;
    DiscordID: string;
    ColorOne: string;
    ColorTwo: string;
    ColorThree: string;
    OverallGrade: string;
    OffenseGrade: string;
    DefenseGrade: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.TeamName = source["TeamName"];
        this.Mascot = source["Mascot"];
        this.Abbreviation = source["Abbreviation"];
        this.ConferenceID = source["ConferenceID"];
        this.Conference = source["Conference"];
        this.Coach = source["Coach"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.ArenaID = source["ArenaID"];
        this.Arena = source["Arena"];
        this.ArenaCapacity = source["ArenaCapacity"];
        this.RecordAttendance = source["RecordAttendance"];
        this.FirstPlayed = source["FirstPlayed"];
        this.DiscordID = source["DiscordID"];
        this.ColorOne = source["ColorOne"];
        this.ColorTwo = source["ColorTwo"];
        this.ColorThree = source["ColorThree"];
        this.OverallGrade = source["OverallGrade"];
        this.OffenseGrade = source["OffenseGrade"];
        this.DefenseGrade = source["DefenseGrade"];
    }
}
export class CollegeTeam {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamName: string;
    Mascot: string;
    Abbreviation: string;
    ConferenceID: number;
    Conference: string;
    Coach: string;
    City: string;
    State: string;
    Country: string;
    ArenaID: number;
    Arena: string;
    ArenaCapacity: number;
    RecordAttendance: number;
    FirstPlayed: number;
    DiscordID: string;
    ColorOne: string;
    ColorTwo: string;
    ColorThree: string;
    OverallGrade: string;
    OffenseGrade: string;
    DefenseGrade: string;
    IsUserCoached: boolean;
    IsClub: boolean;
    IsActive: boolean;
    PlayersProgressed: boolean;
    ProgramPrestige: number;
    ProfessionalPrestige: number;
    Traditions: number;
    Facilities: number;
    Atmosphere: number;
    Academics: number;
    ConferencePrestige: number;
    CoachRating: number;
    SeasonMomentum: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamName = source["TeamName"];
        this.Mascot = source["Mascot"];
        this.Abbreviation = source["Abbreviation"];
        this.ConferenceID = source["ConferenceID"];
        this.Conference = source["Conference"];
        this.Coach = source["Coach"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.ArenaID = source["ArenaID"];
        this.Arena = source["Arena"];
        this.ArenaCapacity = source["ArenaCapacity"];
        this.RecordAttendance = source["RecordAttendance"];
        this.FirstPlayed = source["FirstPlayed"];
        this.DiscordID = source["DiscordID"];
        this.ColorOne = source["ColorOne"];
        this.ColorTwo = source["ColorTwo"];
        this.ColorThree = source["ColorThree"];
        this.OverallGrade = source["OverallGrade"];
        this.OffenseGrade = source["OffenseGrade"];
        this.DefenseGrade = source["DefenseGrade"];
        this.IsUserCoached = source["IsUserCoached"];
        this.IsClub = source["IsClub"];
        this.IsActive = source["IsActive"];
        this.PlayersProgressed = source["PlayersProgressed"];
        this.ProgramPrestige = source["ProgramPrestige"];
        this.ProfessionalPrestige = source["ProfessionalPrestige"];
        this.Traditions = source["Traditions"];
        this.Facilities = source["Facilities"];
        this.Atmosphere = source["Atmosphere"];
        this.Academics = source["Academics"];
        this.ConferencePrestige = source["ConferencePrestige"];
        this.CoachRating = source["CoachRating"];
        this.SeasonMomentum = source["SeasonMomentum"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalTeam {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamName: string;
    Mascot: string;
    Abbreviation: string;
    ConferenceID: number;
    Conference: string;
    Coach: string;
    City: string;
    State: string;
    Country: string;
    ArenaID: number;
    Arena: string;
    ArenaCapacity: number;
    RecordAttendance: number;
    FirstPlayed: number;
    DiscordID: string;
    ColorOne: string;
    ColorTwo: string;
    ColorThree: string;
    OverallGrade: string;
    OffenseGrade: string;
    DefenseGrade: string;
    Owner: string;
    GM: string;
    Scout: string;
    Marketing: string;
    DivisionID: number;
    Division: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamName = source["TeamName"];
        this.Mascot = source["Mascot"];
        this.Abbreviation = source["Abbreviation"];
        this.ConferenceID = source["ConferenceID"];
        this.Conference = source["Conference"];
        this.Coach = source["Coach"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.ArenaID = source["ArenaID"];
        this.Arena = source["Arena"];
        this.ArenaCapacity = source["ArenaCapacity"];
        this.RecordAttendance = source["RecordAttendance"];
        this.FirstPlayed = source["FirstPlayed"];
        this.DiscordID = source["DiscordID"];
        this.ColorOne = source["ColorOne"];
        this.ColorTwo = source["ColorTwo"];
        this.ColorThree = source["ColorThree"];
        this.OverallGrade = source["OverallGrade"];
        this.OffenseGrade = source["OffenseGrade"];
        this.DefenseGrade = source["DefenseGrade"];
        this.Owner = source["Owner"];
        this.GM = source["GM"];
        this.Scout = source["Scout"];
        this.Marketing = source["Marketing"];
        this.DivisionID = source["DivisionID"];
        this.Division = source["Division"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalTeamFranchise {
    TeamID: number;
    HomeMarket: number;
    TicketPrice: number;
    TicketBoxPrice: number;
    TicketDemand: number;
    BoxDemand: number;
    Food1ID: number;
    Food1Price: number;
    Food2ID: number;
    Food2Price: number;
    Food3ID: number;
    Food3Price: number;
    Drink1ID: number;
    Drink1Price: number;
    Drink2ID: number;
    Drink2Price: number;
    Drink3ID: number;
    Drink3Price: number;
    Drink4ID: number;
    Drink4Price: number;
    TeamShirtPrice: number;
    TeamHatPrice: number;
    TeamBeaniePrice: number;
    TeamJerseyPrice: number;
    FacilitiesMaintenanceCost: number;
    EquipmentCost: number;
    OperationsCost: number;
    BathroomPrice: number;
    FanHappinessRating: number;
    AtmosphereRating: number;
    FoodRating: number;
    EmployeeHappinessRating: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.TeamID = source["TeamID"];
        this.HomeMarket = source["HomeMarket"];
        this.TicketPrice = source["TicketPrice"];
        this.TicketBoxPrice = source["TicketBoxPrice"];
        this.TicketDemand = source["TicketDemand"];
        this.BoxDemand = source["BoxDemand"];
        this.Food1ID = source["Food1ID"];
        this.Food1Price = source["Food1Price"];
        this.Food2ID = source["Food2ID"];
        this.Food2Price = source["Food2Price"];
        this.Food3ID = source["Food3ID"];
        this.Food3Price = source["Food3Price"];
        this.Drink1ID = source["Drink1ID"];
        this.Drink1Price = source["Drink1Price"];
        this.Drink2ID = source["Drink2ID"];
        this.Drink2Price = source["Drink2Price"];
        this.Drink3ID = source["Drink3ID"];
        this.Drink3Price = source["Drink3Price"];
        this.Drink4ID = source["Drink4ID"];
        this.Drink4Price = source["Drink4Price"];
        this.TeamShirtPrice = source["TeamShirtPrice"];
        this.TeamHatPrice = source["TeamHatPrice"];
        this.TeamBeaniePrice = source["TeamBeaniePrice"];
        this.TeamJerseyPrice = source["TeamJerseyPrice"];
        this.FacilitiesMaintenanceCost = source["FacilitiesMaintenanceCost"];
        this.EquipmentCost = source["EquipmentCost"];
        this.OperationsCost = source["OperationsCost"];
        this.BathroomPrice = source["BathroomPrice"];
        this.FanHappinessRating = source["FanHappinessRating"];
        this.AtmosphereRating = source["AtmosphereRating"];
        this.FoodRating = source["FoodRating"];
        this.EmployeeHappinessRating = source["EmployeeHappinessRating"];
    }
}
export class Arena {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    Name: string;
    TeamID: number;
    City: string;
    State: string;
    Country: string;
    Capacity: number;
    RecordAttendance: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.Name = source["Name"];
        this.TeamID = source["TeamID"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.Capacity = source["Capacity"];
        this.RecordAttendance = source["RecordAttendance"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class FaceDataResponse {
    PlayerID: number;
    Accessories: string;
    Body: string;
    Ear: string;
    Eye: string;
    EyeLine: string;
    Eyebrow: string;
    FacialHair: string;
    Glasses: string;
    Hair: string;
    HairBG: string;
    HairFlip: boolean;
    Head: string;
    Jersey: string;
    MiscLine: string;
    Mouth: string;
    MouthFlip: boolean;
    Nose: string;
    NoseFlip: boolean;
    SmileLine: string;
    BodySize: number;
    EarSize: number;
    EyeAngle: number;
    EyeBrowAngle: number;
    FaceSize: number;
    FacialHairShave: string;
    NoseSize: number;
    SmileLineSize: number;
    SkinColor: string;
    HairColor: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.PlayerID = source["PlayerID"];
        this.Accessories = source["Accessories"];
        this.Body = source["Body"];
        this.Ear = source["Ear"];
        this.Eye = source["Eye"];
        this.EyeLine = source["EyeLine"];
        this.Eyebrow = source["Eyebrow"];
        this.FacialHair = source["FacialHair"];
        this.Glasses = source["Glasses"];
        this.Hair = source["Hair"];
        this.HairBG = source["HairBG"];
        this.HairFlip = source["HairFlip"];
        this.Head = source["Head"];
        this.Jersey = source["Jersey"];
        this.MiscLine = source["MiscLine"];
        this.Mouth = source["Mouth"];
        this.MouthFlip = source["MouthFlip"];
        this.Nose = source["Nose"];
        this.NoseFlip = source["NoseFlip"];
        this.SmileLine = source["SmileLine"];
        this.BodySize = source["BodySize"];
        this.EarSize = source["EarSize"];
        this.EyeAngle = source["EyeAngle"];
        this.EyeBrowAngle = source["EyeBrowAngle"];
        this.FaceSize = source["FaceSize"];
        this.FacialHairShave = source["FacialHairShave"];
        this.NoseSize = source["NoseSize"];
        this.SmileLineSize = source["SmileLineSize"];
        this.SkinColor = source["SkinColor"];
        this.HairColor = source["HairColor"];
    }
}
export class ProfessionalShootoutLineup {
    ID: number;
    TeamID: number;
    Shooter1ID: number;
    Shooter1ShotType: number;
    Shooter2ID: number;
    Shooter2ShotType: number;
    Shooter3ID: number;
    Shooter3ShotType: number;
    Shooter4ID: number;
    Shooter4ShotType: number;
    Shooter5ID: number;
    Shooter5ShotType: number;
    Shooter6ID: number;
    Shooter6ShotType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.TeamID = source["TeamID"];
        this.Shooter1ID = source["Shooter1ID"];
        this.Shooter1ShotType = source["Shooter1ShotType"];
        this.Shooter2ID = source["Shooter2ID"];
        this.Shooter2ShotType = source["Shooter2ShotType"];
        this.Shooter3ID = source["Shooter3ID"];
        this.Shooter3ShotType = source["Shooter3ShotType"];
        this.Shooter4ID = source["Shooter4ID"];
        this.Shooter4ShotType = source["Shooter4ShotType"];
        this.Shooter5ID = source["Shooter5ID"];
        this.Shooter5ShotType = source["Shooter5ShotType"];
        this.Shooter6ID = source["Shooter6ID"];
        this.Shooter6ShotType = source["Shooter6ShotType"];
    }
}
export class ProfessionalLineup {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    Line: number;
    LineType: number;
    CenterID: number;
    Forward1ID: number;
    Forward2ID: number;
    Defender1ID: number;
    Defender2ID: number;
    GoalieID: number;
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.Line = source["Line"];
        this.LineType = source["LineType"];
        this.CenterID = source["CenterID"];
        this.Forward1ID = source["Forward1ID"];
        this.Forward2ID = source["Forward2ID"];
        this.Defender1ID = source["Defender1ID"];
        this.Defender2ID = source["Defender2ID"];
        this.GoalieID = source["GoalieID"];
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalGame {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    WeekID: number;
    Week: number;
    SeasonID: number;
    GameDay: string;
    HomeTeamRank: number;
    HomeTeamID: number;
    HomeTeam: string;
    HomeTeamCoach: string;
    HomeTeamWin: boolean;
    AwayTeamRank: number;
    AwayTeamID: number;
    AwayTeam: string;
    AwayTeamCoach: string;
    AwayTeamWin: boolean;
    StarOne: number;
    StarTwo: number;
    StarThree: number;
    HomeTeamScore: number;
    AwayTeamScore: number;
    HomeTeamShootoutScore: number;
    AwayTeamShootoutScore: number;
    ArenaID: number;
    Arena: string;
    AttendanceCount: number;
    City: string;
    State: string;
    Country: string;
    IsNeutralSite: boolean;
    IsConference: boolean;
    IsPlayoffGame: boolean;
    IsRivalryGame: boolean;
    GameComplete: boolean;
    IsOvertime: boolean;
    IsShootout: boolean;
    GameTitle: string;
    NextGameID: number;
    NextGameHOA: string;
    IsPreseason: boolean;
    SeriesID: number;
    IsDivisional: boolean;
    IsStanleyCup: boolean;
    IsInternational: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.WeekID = source["WeekID"];
        this.Week = source["Week"];
        this.SeasonID = source["SeasonID"];
        this.GameDay = source["GameDay"];
        this.HomeTeamRank = source["HomeTeamRank"];
        this.HomeTeamID = source["HomeTeamID"];
        this.HomeTeam = source["HomeTeam"];
        this.HomeTeamCoach = source["HomeTeamCoach"];
        this.HomeTeamWin = source["HomeTeamWin"];
        this.AwayTeamRank = source["AwayTeamRank"];
        this.AwayTeamID = source["AwayTeamID"];
        this.AwayTeam = source["AwayTeam"];
        this.AwayTeamCoach = source["AwayTeamCoach"];
        this.AwayTeamWin = source["AwayTeamWin"];
        this.StarOne = source["StarOne"];
        this.StarTwo = source["StarTwo"];
        this.StarThree = source["StarThree"];
        this.HomeTeamScore = source["HomeTeamScore"];
        this.AwayTeamScore = source["AwayTeamScore"];
        this.HomeTeamShootoutScore = source["HomeTeamShootoutScore"];
        this.AwayTeamShootoutScore = source["AwayTeamShootoutScore"];
        this.ArenaID = source["ArenaID"];
        this.Arena = source["Arena"];
        this.AttendanceCount = source["AttendanceCount"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.IsNeutralSite = source["IsNeutralSite"];
        this.IsConference = source["IsConference"];
        this.IsPlayoffGame = source["IsPlayoffGame"];
        this.IsRivalryGame = source["IsRivalryGame"];
        this.GameComplete = source["GameComplete"];
        this.IsOvertime = source["IsOvertime"];
        this.IsShootout = source["IsShootout"];
        this.GameTitle = source["GameTitle"];
        this.NextGameID = source["NextGameID"];
        this.NextGameHOA = source["NextGameHOA"];
        this.IsPreseason = source["IsPreseason"];
        this.SeriesID = source["SeriesID"];
        this.IsDivisional = source["IsDivisional"];
        this.IsStanleyCup = source["IsStanleyCup"];
        this.IsInternational = source["IsInternational"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProCapsheet {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    Y1Salary: number;
    Y1CapHit: number;
    Y2Salary: number;
    Y2CapHit: number;
    Y3Salary: number;
    Y3CapHit: number;
    Y4Salary: number;
    Y4CapHit: number;
    Y5Salary: number;
    Y5CapHit: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.Y1Salary = source["Y1Salary"];
        this.Y1CapHit = source["Y1CapHit"];
        this.Y2Salary = source["Y2Salary"];
        this.Y2CapHit = source["Y2CapHit"];
        this.Y3Salary = source["Y3Salary"];
        this.Y3CapHit = source["Y3CapHit"];
        this.Y4Salary = source["Y4Salary"];
        this.Y4CapHit = source["Y4CapHit"];
        this.Y5Salary = source["Y5Salary"];
        this.Y5CapHit = source["Y5CapHit"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProfessionalStandings {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    TeamName: string;
    SeasonID: number;
    Season: number;
    LeagueID: number;
    ConferenceID: number;
    TotalWins: number;
    TotalLosses: number;
    TotalOTWins: number;
    TotalOTLosses: number;
    ShootoutWins: number;
    ShootoutLosses: number;
    ConferenceWins: number;
    ConferenceLosses: number;
    ConferenceOTWins: number;
    ConferenceOTLosses: number;
    RankedWins: number;
    RankedLosses: number;
    Points: number;
    GoalsFor: number;
    GoalsAgainst: number;
    Streak: number;
    IsWinStreak: boolean;
    HomeWins: number;
    AwayWins: number;
    PostSeasonStatus: string;
    DivisionID: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.TeamName = source["TeamName"];
        this.SeasonID = source["SeasonID"];
        this.Season = source["Season"];
        this.LeagueID = source["LeagueID"];
        this.ConferenceID = source["ConferenceID"];
        this.TotalWins = source["TotalWins"];
        this.TotalLosses = source["TotalLosses"];
        this.TotalOTWins = source["TotalOTWins"];
        this.TotalOTLosses = source["TotalOTLosses"];
        this.ShootoutWins = source["ShootoutWins"];
        this.ShootoutLosses = source["ShootoutLosses"];
        this.ConferenceWins = source["ConferenceWins"];
        this.ConferenceLosses = source["ConferenceLosses"];
        this.ConferenceOTWins = source["ConferenceOTWins"];
        this.ConferenceOTLosses = source["ConferenceOTLosses"];
        this.RankedWins = source["RankedWins"];
        this.RankedLosses = source["RankedLosses"];
        this.Points = source["Points"];
        this.GoalsFor = source["GoalsFor"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.Streak = source["Streak"];
        this.IsWinStreak = source["IsWinStreak"];
        this.HomeWins = source["HomeWins"];
        this.AwayWins = source["AwayWins"];
        this.PostSeasonStatus = source["PostSeasonStatus"];
        this.DivisionID = source["DivisionID"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegePollOfficial {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    WeekID: number;
    Week: number;
    Rank1: string;
    Rank1ID: number;
    Rank1Votes: number;
    Rank1No1Votes: number;
    Rank2: string;
    Rank2ID: number;
    Rank2Votes: number;
    Rank2No1Votes: number;
    Rank3: string;
    Rank3ID: number;
    Rank3Votes: number;
    Rank3No1Votes: number;
    Rank4: string;
    Rank4ID: number;
    Rank4Votes: number;
    Rank4No1Votes: number;
    Rank5: string;
    Rank5ID: number;
    Rank5Votes: number;
    Rank5No1Votes: number;
    Rank6: string;
    Rank6ID: number;
    Rank6Votes: number;
    Rank6No1Votes: number;
    Rank7: string;
    Rank7ID: number;
    Rank7Votes: number;
    Rank7No1Votes: number;
    Rank8: string;
    Rank8ID: number;
    Rank8Votes: number;
    Rank8No1Votes: number;
    Rank9: string;
    Rank9ID: number;
    Rank9Votes: number;
    Rank9No1Votes: number;
    Rank10: string;
    Rank10ID: number;
    Rank10Votes: number;
    Rank10No1Votes: number;
    Rank11: string;
    Rank11ID: number;
    Rank11Votes: number;
    Rank11No1Votes: number;
    Rank12: string;
    Rank12ID: number;
    Rank12Votes: number;
    Rank12No1Votes: number;
    Rank13: string;
    Rank13ID: number;
    Rank13Votes: number;
    Rank13No1Votes: number;
    Rank14: string;
    Rank14ID: number;
    Rank14Votes: number;
    Rank14No1Votes: number;
    Rank15: string;
    Rank15ID: number;
    Rank15Votes: number;
    Rank15No1Votes: number;
    Rank16: string;
    Rank16ID: number;
    Rank16Votes: number;
    Rank16No1Votes: number;
    Rank17: string;
    Rank17ID: number;
    Rank17Votes: number;
    Rank17No1Votes: number;
    Rank18: string;
    Rank18ID: number;
    Rank18Votes: number;
    Rank18No1Votes: number;
    Rank19: string;
    Rank19ID: number;
    Rank19Votes: number;
    Rank19No1Votes: number;
    Rank20: string;
    Rank20ID: number;
    Rank20Votes: number;
    Rank20No1Votes: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.WeekID = source["WeekID"];
        this.Week = source["Week"];
        this.Rank1 = source["Rank1"];
        this.Rank1ID = source["Rank1ID"];
        this.Rank1Votes = source["Rank1Votes"];
        this.Rank1No1Votes = source["Rank1No1Votes"];
        this.Rank2 = source["Rank2"];
        this.Rank2ID = source["Rank2ID"];
        this.Rank2Votes = source["Rank2Votes"];
        this.Rank2No1Votes = source["Rank2No1Votes"];
        this.Rank3 = source["Rank3"];
        this.Rank3ID = source["Rank3ID"];
        this.Rank3Votes = source["Rank3Votes"];
        this.Rank3No1Votes = source["Rank3No1Votes"];
        this.Rank4 = source["Rank4"];
        this.Rank4ID = source["Rank4ID"];
        this.Rank4Votes = source["Rank4Votes"];
        this.Rank4No1Votes = source["Rank4No1Votes"];
        this.Rank5 = source["Rank5"];
        this.Rank5ID = source["Rank5ID"];
        this.Rank5Votes = source["Rank5Votes"];
        this.Rank5No1Votes = source["Rank5No1Votes"];
        this.Rank6 = source["Rank6"];
        this.Rank6ID = source["Rank6ID"];
        this.Rank6Votes = source["Rank6Votes"];
        this.Rank6No1Votes = source["Rank6No1Votes"];
        this.Rank7 = source["Rank7"];
        this.Rank7ID = source["Rank7ID"];
        this.Rank7Votes = source["Rank7Votes"];
        this.Rank7No1Votes = source["Rank7No1Votes"];
        this.Rank8 = source["Rank8"];
        this.Rank8ID = source["Rank8ID"];
        this.Rank8Votes = source["Rank8Votes"];
        this.Rank8No1Votes = source["Rank8No1Votes"];
        this.Rank9 = source["Rank9"];
        this.Rank9ID = source["Rank9ID"];
        this.Rank9Votes = source["Rank9Votes"];
        this.Rank9No1Votes = source["Rank9No1Votes"];
        this.Rank10 = source["Rank10"];
        this.Rank10ID = source["Rank10ID"];
        this.Rank10Votes = source["Rank10Votes"];
        this.Rank10No1Votes = source["Rank10No1Votes"];
        this.Rank11 = source["Rank11"];
        this.Rank11ID = source["Rank11ID"];
        this.Rank11Votes = source["Rank11Votes"];
        this.Rank11No1Votes = source["Rank11No1Votes"];
        this.Rank12 = source["Rank12"];
        this.Rank12ID = source["Rank12ID"];
        this.Rank12Votes = source["Rank12Votes"];
        this.Rank12No1Votes = source["Rank12No1Votes"];
        this.Rank13 = source["Rank13"];
        this.Rank13ID = source["Rank13ID"];
        this.Rank13Votes = source["Rank13Votes"];
        this.Rank13No1Votes = source["Rank13No1Votes"];
        this.Rank14 = source["Rank14"];
        this.Rank14ID = source["Rank14ID"];
        this.Rank14Votes = source["Rank14Votes"];
        this.Rank14No1Votes = source["Rank14No1Votes"];
        this.Rank15 = source["Rank15"];
        this.Rank15ID = source["Rank15ID"];
        this.Rank15Votes = source["Rank15Votes"];
        this.Rank15No1Votes = source["Rank15No1Votes"];
        this.Rank16 = source["Rank16"];
        this.Rank16ID = source["Rank16ID"];
        this.Rank16Votes = source["Rank16Votes"];
        this.Rank16No1Votes = source["Rank16No1Votes"];
        this.Rank17 = source["Rank17"];
        this.Rank17ID = source["Rank17ID"];
        this.Rank17Votes = source["Rank17Votes"];
        this.Rank17No1Votes = source["Rank17No1Votes"];
        this.Rank18 = source["Rank18"];
        this.Rank18ID = source["Rank18ID"];
        this.Rank18Votes = source["Rank18Votes"];
        this.Rank18No1Votes = source["Rank18No1Votes"];
        this.Rank19 = source["Rank19"];
        this.Rank19ID = source["Rank19ID"];
        this.Rank19Votes = source["Rank19Votes"];
        this.Rank19No1Votes = source["Rank19No1Votes"];
        this.Rank20 = source["Rank20"];
        this.Rank20ID = source["Rank20ID"];
        this.Rank20Votes = source["Rank20Votes"];
        this.Rank20No1Votes = source["Rank20No1Votes"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegePollSubmission {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    Username: string;
    SeasonID: number;
    WeekID: number;
    Week: number;
    Rank1: string;
    Rank1ID: number;
    Rank2: string;
    Rank2ID: number;
    Rank3: string;
    Rank3ID: number;
    Rank4: string;
    Rank4ID: number;
    Rank5: string;
    Rank5ID: number;
    Rank6: string;
    Rank6ID: number;
    Rank7: string;
    Rank7ID: number;
    Rank8: string;
    Rank8ID: number;
    Rank9: string;
    Rank9ID: number;
    Rank10: string;
    Rank10ID: number;
    Rank11: string;
    Rank11ID: number;
    Rank12: string;
    Rank12ID: number;
    Rank13: string;
    Rank13ID: number;
    Rank14: string;
    Rank14ID: number;
    Rank15: string;
    Rank15ID: number;
    Rank16: string;
    Rank16ID: number;
    Rank17: string;
    Rank17ID: number;
    Rank18: string;
    Rank18ID: number;
    Rank19: string;
    Rank19ID: number;
    Rank20: string;
    Rank20ID: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.Username = source["Username"];
        this.SeasonID = source["SeasonID"];
        this.WeekID = source["WeekID"];
        this.Week = source["Week"];
        this.Rank1 = source["Rank1"];
        this.Rank1ID = source["Rank1ID"];
        this.Rank2 = source["Rank2"];
        this.Rank2ID = source["Rank2ID"];
        this.Rank3 = source["Rank3"];
        this.Rank3ID = source["Rank3ID"];
        this.Rank4 = source["Rank4"];
        this.Rank4ID = source["Rank4ID"];
        this.Rank5 = source["Rank5"];
        this.Rank5ID = source["Rank5ID"];
        this.Rank6 = source["Rank6"];
        this.Rank6ID = source["Rank6ID"];
        this.Rank7 = source["Rank7"];
        this.Rank7ID = source["Rank7ID"];
        this.Rank8 = source["Rank8"];
        this.Rank8ID = source["Rank8ID"];
        this.Rank9 = source["Rank9"];
        this.Rank9ID = source["Rank9ID"];
        this.Rank10 = source["Rank10"];
        this.Rank10ID = source["Rank10ID"];
        this.Rank11 = source["Rank11"];
        this.Rank11ID = source["Rank11ID"];
        this.Rank12 = source["Rank12"];
        this.Rank12ID = source["Rank12ID"];
        this.Rank13 = source["Rank13"];
        this.Rank13ID = source["Rank13ID"];
        this.Rank14 = source["Rank14"];
        this.Rank14ID = source["Rank14ID"];
        this.Rank15 = source["Rank15"];
        this.Rank15ID = source["Rank15ID"];
        this.Rank16 = source["Rank16"];
        this.Rank16ID = source["Rank16ID"];
        this.Rank17 = source["Rank17"];
        this.Rank17ID = source["Rank17ID"];
        this.Rank18 = source["Rank18"];
        this.Rank18ID = source["Rank18ID"];
        this.Rank19 = source["Rank19"];
        this.Rank19ID = source["Rank19ID"];
        this.Rank20 = source["Rank20"];
        this.Rank20ID = source["Rank20ID"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegeShootoutLineup {
    ID: number;
    TeamID: number;
    Shooter1ID: number;
    Shooter1ShotType: number;
    Shooter2ID: number;
    Shooter2ShotType: number;
    Shooter3ID: number;
    Shooter3ShotType: number;
    Shooter4ID: number;
    Shooter4ShotType: number;
    Shooter5ID: number;
    Shooter5ShotType: number;
    Shooter6ID: number;
    Shooter6ShotType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.TeamID = source["TeamID"];
        this.Shooter1ID = source["Shooter1ID"];
        this.Shooter1ShotType = source["Shooter1ShotType"];
        this.Shooter2ID = source["Shooter2ID"];
        this.Shooter2ShotType = source["Shooter2ShotType"];
        this.Shooter3ID = source["Shooter3ID"];
        this.Shooter3ShotType = source["Shooter3ShotType"];
        this.Shooter4ID = source["Shooter4ID"];
        this.Shooter4ShotType = source["Shooter4ShotType"];
        this.Shooter5ID = source["Shooter5ID"];
        this.Shooter5ShotType = source["Shooter5ShotType"];
        this.Shooter6ID = source["Shooter6ID"];
        this.Shooter6ShotType = source["Shooter6ShotType"];
    }
}
export class CollegeLineup {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    Line: number;
    LineType: number;
    CenterID: number;
    Forward1ID: number;
    Forward2ID: number;
    Defender1ID: number;
    Defender2ID: number;
    GoalieID: number;
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.Line = source["Line"];
        this.LineType = source["LineType"];
        this.CenterID = source["CenterID"];
        this.Forward1ID = source["Forward1ID"];
        this.Forward2ID = source["Forward2ID"];
        this.Defender1ID = source["Defender1ID"];
        this.Defender2ID = source["Defender2ID"];
        this.GoalieID = source["GoalieID"];
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegeGame {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    WeekID: number;
    Week: number;
    SeasonID: number;
    GameDay: string;
    HomeTeamRank: number;
    HomeTeamID: number;
    HomeTeam: string;
    HomeTeamCoach: string;
    HomeTeamWin: boolean;
    AwayTeamRank: number;
    AwayTeamID: number;
    AwayTeam: string;
    AwayTeamCoach: string;
    AwayTeamWin: boolean;
    StarOne: number;
    StarTwo: number;
    StarThree: number;
    HomeTeamScore: number;
    AwayTeamScore: number;
    HomeTeamShootoutScore: number;
    AwayTeamShootoutScore: number;
    ArenaID: number;
    Arena: string;
    AttendanceCount: number;
    City: string;
    State: string;
    Country: string;
    IsNeutralSite: boolean;
    IsConference: boolean;
    IsPlayoffGame: boolean;
    IsRivalryGame: boolean;
    GameComplete: boolean;
    IsOvertime: boolean;
    IsShootout: boolean;
    GameTitle: string;
    NextGameID: number;
    NextGameHOA: string;
    IsPreseason: boolean;
    IsNationalChampionship: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.WeekID = source["WeekID"];
        this.Week = source["Week"];
        this.SeasonID = source["SeasonID"];
        this.GameDay = source["GameDay"];
        this.HomeTeamRank = source["HomeTeamRank"];
        this.HomeTeamID = source["HomeTeamID"];
        this.HomeTeam = source["HomeTeam"];
        this.HomeTeamCoach = source["HomeTeamCoach"];
        this.HomeTeamWin = source["HomeTeamWin"];
        this.AwayTeamRank = source["AwayTeamRank"];
        this.AwayTeamID = source["AwayTeamID"];
        this.AwayTeam = source["AwayTeam"];
        this.AwayTeamCoach = source["AwayTeamCoach"];
        this.AwayTeamWin = source["AwayTeamWin"];
        this.StarOne = source["StarOne"];
        this.StarTwo = source["StarTwo"];
        this.StarThree = source["StarThree"];
        this.HomeTeamScore = source["HomeTeamScore"];
        this.AwayTeamScore = source["AwayTeamScore"];
        this.HomeTeamShootoutScore = source["HomeTeamShootoutScore"];
        this.AwayTeamShootoutScore = source["AwayTeamShootoutScore"];
        this.ArenaID = source["ArenaID"];
        this.Arena = source["Arena"];
        this.AttendanceCount = source["AttendanceCount"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.IsNeutralSite = source["IsNeutralSite"];
        this.IsConference = source["IsConference"];
        this.IsPlayoffGame = source["IsPlayoffGame"];
        this.IsRivalryGame = source["IsRivalryGame"];
        this.GameComplete = source["GameComplete"];
        this.IsOvertime = source["IsOvertime"];
        this.IsShootout = source["IsShootout"];
        this.GameTitle = source["GameTitle"];
        this.NextGameID = source["NextGameID"];
        this.NextGameHOA = source["NextGameHOA"];
        this.IsPreseason = source["IsPreseason"];
        this.IsNationalChampionship = source["IsNationalChampionship"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Notification {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    League: string;
    NotificationType: string;
    Message: string;
    Subject: string;
    IsRead: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.League = source["League"];
        this.NotificationType = source["NotificationType"];
        this.Message = source["Message"];
        this.Subject = source["Subject"];
        this.IsRead = source["IsRead"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class NewsLog {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    WeekID: number;
    Week: number;
    SeasonID: number;
    Season: number;
    MessageType: string;
    Message: string;
    League: string;
    TeamID: number;
    ShowLog: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.WeekID = source["WeekID"];
        this.Week = source["Week"];
        this.SeasonID = source["SeasonID"];
        this.Season = source["Season"];
        this.MessageType = source["MessageType"];
        this.Message = source["Message"];
        this.League = source["League"];
        this.TeamID = source["TeamID"];
        this.ShowLog = source["ShowLog"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class RecruitingTeamProfile {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    Team: string;
    State: string;
    Country: string;
    ScholarshipsAvailable: number;
    WeeklyPoints: number;
    WeeklyScoutingPoints: number;
    SpentPoints: number;
    TotalCommitments: number;
    RecruitClassSize: number;
    PortalReputation: number;
    ESPNScore: number;
    RivalsScore: number;
    Rank247Score: number;
    CompositeScore: number;
    ThreeStars: number;
    FourStars: number;
    FiveStars: number;
    RecruitingClassRank: number;
    CaughtCheating: boolean;
    IsAI: boolean;
    IsUserTeam: boolean;
    AIBehavior: string;
    AIQuality: string;
    WeeksMissed: number;
    BattlesWon: number;
    BattlesLost: number;
    AIMinThreshold: number;
    AIMaxThreshold: number;
    AIStarMin: number;
    AIStarMax: number;
    Recruiter: string;
    Recruits: RecruitPlayerProfile[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.ScholarshipsAvailable = source["ScholarshipsAvailable"];
        this.WeeklyPoints = source["WeeklyPoints"];
        this.WeeklyScoutingPoints = source["WeeklyScoutingPoints"];
        this.SpentPoints = source["SpentPoints"];
        this.TotalCommitments = source["TotalCommitments"];
        this.RecruitClassSize = source["RecruitClassSize"];
        this.PortalReputation = source["PortalReputation"];
        this.ESPNScore = source["ESPNScore"];
        this.RivalsScore = source["RivalsScore"];
        this.Rank247Score = source["Rank247Score"];
        this.CompositeScore = source["CompositeScore"];
        this.ThreeStars = source["ThreeStars"];
        this.FourStars = source["FourStars"];
        this.FiveStars = source["FiveStars"];
        this.RecruitingClassRank = source["RecruitingClassRank"];
        this.CaughtCheating = source["CaughtCheating"];
        this.IsAI = source["IsAI"];
        this.IsUserTeam = source["IsUserTeam"];
        this.AIBehavior = source["AIBehavior"];
        this.AIQuality = source["AIQuality"];
        this.WeeksMissed = source["WeeksMissed"];
        this.BattlesWon = source["BattlesWon"];
        this.BattlesLost = source["BattlesLost"];
        this.AIMinThreshold = source["AIMinThreshold"];
        this.AIMaxThreshold = source["AIMaxThreshold"];
        this.AIStarMin = source["AIStarMin"];
        this.AIStarMax = source["AIStarMax"];
        this.Recruiter = source["Recruiter"];
        this.Recruits = this.convertValues(source["Recruits"], RecruitPlayerProfile);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Recruit {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    FirstName: string;
    LastName: string;
    Position: string;
    Archetype: string;
    TeamID: number;
    Team: string;
    Height: number;
    Weight: number;
    Stars: number;
    Age: number;
    Overall: number;
    Agility: number;
    Faceoffs: number;
    LongShotAccuracy: number;
    LongShotPower: number;
    CloseShotAccuracy: number;
    CloseShotPower: number;
    OneTimer: number;
    Passing: number;
    PuckHandling: number;
    Strength: number;
    BodyChecking: number;
    StickChecking: number;
    ShotBlocking: number;
    Goalkeeping: number;
    GoalieVision: number;
    GoalieReboundControl: number;
    Discipline: number;
    Aggression: number;
    Stamina: number;
    GoalieStamina: number;
    InjuryRating: number;
    DisciplineDeviation: number;
    InjuryDeviation: number;
    PrimeAge: number;
    Clutch: number;
    HighSchool: string;
    City: string;
    State: string;
    Country: string;
    OriginalTeamID: number;
    OriginalTeam: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Competitiveness: number;
    TeamLoyalty: number;
    PlaytimePreference: number;
    PlayerMorale: number;
    Personality: string;
    RelativeID: number;
    RelativeType: number;
    Notes: string;
    HasProgressed: boolean;
    ProgramPref: number;
    ProfDevPref: number;
    TraditionsPref: number;
    FacilitiesPref: number;
    AtmospherePref: number;
    AcademicsPref: number;
    ConferencePref: number;
    CoachPref: number;
    SeasonMomentumPref: number;
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;
    AgilityPotential: number;
    FaceoffsPotential: number;
    CloseShotAccuracyPotential: number;
    CloseShotPowerPotential: number;
    LongShotAccuracyPotential: number;
    LongShotPowerPotential: number;
    PassingPotential: number;
    PuckHandlingPotential: number;
    StrengthPotential: number;
    BodyCheckingPotential: number;
    StickCheckingPotential: number;
    ShotBlockingPotential: number;
    GoalkeepingPotential: number;
    GoalieVisionPotential: number;
    GoalieReboundPotential: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    InjuryCount: number;
    Regression: number;
    DecayRate: number;
    IsSigned: boolean;
    College: string;
    IsCustomCroot: boolean;
    CustomCrootFor: string;
    RecruitModifier: number;
    CompositeRank: number;
    RivalsRank: number;
    ESPNRank: number;
    Rank247: number;
    TopRankModifier: number;
    RecruitingModifier: number;
    RecruitingStatus: string;
    RecruitPlayerProfiles: RecruitPlayerProfile[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.FirstName = source["FirstName"];
        this.LastName = source["LastName"];
        this.Position = source["Position"];
        this.Archetype = source["Archetype"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.Height = source["Height"];
        this.Weight = source["Weight"];
        this.Stars = source["Stars"];
        this.Age = source["Age"];
        this.Overall = source["Overall"];
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.GoalieReboundControl = source["GoalieReboundControl"];
        this.Discipline = source["Discipline"];
        this.Aggression = source["Aggression"];
        this.Stamina = source["Stamina"];
        this.GoalieStamina = source["GoalieStamina"];
        this.InjuryRating = source["InjuryRating"];
        this.DisciplineDeviation = source["DisciplineDeviation"];
        this.InjuryDeviation = source["InjuryDeviation"];
        this.PrimeAge = source["PrimeAge"];
        this.Clutch = source["Clutch"];
        this.HighSchool = source["HighSchool"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.OriginalTeamID = source["OriginalTeamID"];
        this.OriginalTeam = source["OriginalTeam"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Competitiveness = source["Competitiveness"];
        this.TeamLoyalty = source["TeamLoyalty"];
        this.PlaytimePreference = source["PlaytimePreference"];
        this.PlayerMorale = source["PlayerMorale"];
        this.Personality = source["Personality"];
        this.RelativeID = source["RelativeID"];
        this.RelativeType = source["RelativeType"];
        this.Notes = source["Notes"];
        this.HasProgressed = source["HasProgressed"];
        this.ProgramPref = source["ProgramPref"];
        this.ProfDevPref = source["ProfDevPref"];
        this.TraditionsPref = source["TraditionsPref"];
        this.FacilitiesPref = source["FacilitiesPref"];
        this.AtmospherePref = source["AtmospherePref"];
        this.AcademicsPref = source["AcademicsPref"];
        this.ConferencePref = source["ConferencePref"];
        this.CoachPref = source["CoachPref"];
        this.SeasonMomentumPref = source["SeasonMomentumPref"];
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
        this.AgilityPotential = source["AgilityPotential"];
        this.FaceoffsPotential = source["FaceoffsPotential"];
        this.CloseShotAccuracyPotential = source["CloseShotAccuracyPotential"];
        this.CloseShotPowerPotential = source["CloseShotPowerPotential"];
        this.LongShotAccuracyPotential = source["LongShotAccuracyPotential"];
        this.LongShotPowerPotential = source["LongShotPowerPotential"];
        this.PassingPotential = source["PassingPotential"];
        this.PuckHandlingPotential = source["PuckHandlingPotential"];
        this.StrengthPotential = source["StrengthPotential"];
        this.BodyCheckingPotential = source["BodyCheckingPotential"];
        this.StickCheckingPotential = source["StickCheckingPotential"];
        this.ShotBlockingPotential = source["ShotBlockingPotential"];
        this.GoalkeepingPotential = source["GoalkeepingPotential"];
        this.GoalieVisionPotential = source["GoalieVisionPotential"];
        this.GoalieReboundPotential = source["GoalieReboundPotential"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.InjuryCount = source["InjuryCount"];
        this.Regression = source["Regression"];
        this.DecayRate = source["DecayRate"];
        this.IsSigned = source["IsSigned"];
        this.College = source["College"];
        this.IsCustomCroot = source["IsCustomCroot"];
        this.CustomCrootFor = source["CustomCrootFor"];
        this.RecruitModifier = source["RecruitModifier"];
        this.CompositeRank = source["CompositeRank"];
        this.RivalsRank = source["RivalsRank"];
        this.ESPNRank = source["ESPNRank"];
        this.Rank247 = source["Rank247"];
        this.TopRankModifier = source["TopRankModifier"];
        this.RecruitingModifier = source["RecruitingModifier"];
        this.RecruitingStatus = source["RecruitingStatus"];
        this.RecruitPlayerProfiles = this.convertValues(source["RecruitPlayerProfiles"], RecruitPlayerProfile);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class RecruitPlayerProfile {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    RecruitID: number;
    ProfileID: number;
    TotalPoints: number;
    CurrentWeeksPoints: number;
    PreviousWeekPoints: number;
    Modifier: number;
    IsHomeState: boolean;
    IsPipelineState: boolean;
    SpendingCount: number;
    Scholarship: boolean;
    ScholarshipRevoked: boolean;
    RemovedFromBoard: boolean;
    IsSigned: boolean;
    IsLocked: boolean;
    CaughtCheating: boolean;
    TeamReachedMax: boolean;
    Agility: boolean;
    Faceoffs: boolean;
    LongShotAccuracy: boolean;
    LongShotPower: boolean;
    CloseShotAccuracy: boolean;
    CloseShotPower: boolean;
    OneTimer: boolean;
    Passing: boolean;
    PuckHandling: boolean;
    Strength: boolean;
    BodyChecking: boolean;
    StickChecking: boolean;
    ShotBlocking: boolean;
    Goalkeeping: boolean;
    GoalieVision: boolean;
    GoalieReboundControl: boolean;
    Recruit: Recruit;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.RecruitID = source["RecruitID"];
        this.ProfileID = source["ProfileID"];
        this.TotalPoints = source["TotalPoints"];
        this.CurrentWeeksPoints = source["CurrentWeeksPoints"];
        this.PreviousWeekPoints = source["PreviousWeekPoints"];
        this.Modifier = source["Modifier"];
        this.IsHomeState = source["IsHomeState"];
        this.IsPipelineState = source["IsPipelineState"];
        this.SpendingCount = source["SpendingCount"];
        this.Scholarship = source["Scholarship"];
        this.ScholarshipRevoked = source["ScholarshipRevoked"];
        this.RemovedFromBoard = source["RemovedFromBoard"];
        this.IsSigned = source["IsSigned"];
        this.IsLocked = source["IsLocked"];
        this.CaughtCheating = source["CaughtCheating"];
        this.TeamReachedMax = source["TeamReachedMax"];
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.GoalieReboundControl = source["GoalieReboundControl"];
        this.Recruit = this.convertValues(source["Recruit"], Recruit);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class LeadingTeams {
    TeamID: number;
    TeamName: string;
    TeamAbbr: string;
    Odds: number;
    HasScholarship: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.TeamID = source["TeamID"];
        this.TeamName = source["TeamName"];
        this.TeamAbbr = source["TeamAbbr"];
        this.Odds = source["Odds"];
        this.HasScholarship = source["HasScholarship"];
    }
}
export class Croot {
    ID: number;
    TeamID: number;
    College: string;
    FirstName: string;
    LastName: string;
    Position: string;
    Archetype: string;
    Height: number;
    Weight: number;
    Stars: number;
    PotentialGrade: string;
    Personality: string;
    RecruitingBias: string;
    AcademicBias: string;
    WorkEthic: string;
    HighSchool: string;
    City: string;
    State: string;
    Country: string;
    AffinityOne: string;
    AffinityTwo: string;
    RecruitingStatus: string;
    RecruitModifier: number;
    IsCustomCroot: boolean;
    CustomCrootFor: string;
    IsSigned: boolean;
    OverallGrade: string;
    Agility: number;
    Faceoffs: number;
    LongShotAccuracy: number;
    LongShotPower: number;
    CloseShotAccuracy: number;
    CloseShotPower: number;
    OneTimer: number;
    Passing: number;
    PuckHandling: number;
    Strength: number;
    BodyChecking: number;
    StickChecking: number;
    ShotBlocking: number;
    Goalkeeping: number;
    GoalieVision: number;
    TotalRank: number;
    InjuryRating: string;
    Stamina: string;
    AgilityGrade: string;
    FaceoffsGrade: string;
    CloseShotAccuracyGrade: string;
    CloseShotPowerGrade: string;
    LongShotAccuracyGrade: string;
    LongShotPowerGrade: string;
    PassingGrade: string;
    PuckHandlingGrade: string;
    StrengthGrade: string;
    BodyCheckingGrade: string;
    StickCheckingGrade: string;
    ShotBlockingGrade: string;
    GoalkeepingGrade: string;
    GoalieVisionGrade: string;
    GoalieReboundGrade: string;
    ProgramPref: number;
    ProfDevPref: number;
    TraditionsPref: number;
    FacilitiesPref: number;
    AtmospherePref: number;
    AcademicsPref: number;
    ConferencePref: number;
    CoachPref: number;
    SeasonMomentumPref: number;
    LeadingTeams: LeadingTeams[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.TeamID = source["TeamID"];
        this.College = source["College"];
        this.FirstName = source["FirstName"];
        this.LastName = source["LastName"];
        this.Position = source["Position"];
        this.Archetype = source["Archetype"];
        this.Height = source["Height"];
        this.Weight = source["Weight"];
        this.Stars = source["Stars"];
        this.PotentialGrade = source["PotentialGrade"];
        this.Personality = source["Personality"];
        this.RecruitingBias = source["RecruitingBias"];
        this.AcademicBias = source["AcademicBias"];
        this.WorkEthic = source["WorkEthic"];
        this.HighSchool = source["HighSchool"];
        this.City = source["City"];
        this.State = source["State"];
        this.Country = source["Country"];
        this.AffinityOne = source["AffinityOne"];
        this.AffinityTwo = source["AffinityTwo"];
        this.RecruitingStatus = source["RecruitingStatus"];
        this.RecruitModifier = source["RecruitModifier"];
        this.IsCustomCroot = source["IsCustomCroot"];
        this.CustomCrootFor = source["CustomCrootFor"];
        this.IsSigned = source["IsSigned"];
        this.OverallGrade = source["OverallGrade"];
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.TotalRank = source["TotalRank"];
        this.InjuryRating = source["InjuryRating"];
        this.Stamina = source["Stamina"];
        this.AgilityGrade = source["AgilityGrade"];
        this.FaceoffsGrade = source["FaceoffsGrade"];
        this.CloseShotAccuracyGrade = source["CloseShotAccuracyGrade"];
        this.CloseShotPowerGrade = source["CloseShotPowerGrade"];
        this.LongShotAccuracyGrade = source["LongShotAccuracyGrade"];
        this.LongShotPowerGrade = source["LongShotPowerGrade"];
        this.PassingGrade = source["PassingGrade"];
        this.PuckHandlingGrade = source["PuckHandlingGrade"];
        this.StrengthGrade = source["StrengthGrade"];
        this.BodyCheckingGrade = source["BodyCheckingGrade"];
        this.StickCheckingGrade = source["StickCheckingGrade"];
        this.ShotBlockingGrade = source["ShotBlockingGrade"];
        this.GoalkeepingGrade = source["GoalkeepingGrade"];
        this.GoalieVisionGrade = source["GoalieVisionGrade"];
        this.GoalieReboundGrade = source["GoalieReboundGrade"];
        this.ProgramPref = source["ProgramPref"];
        this.ProfDevPref = source["ProfDevPref"];
        this.TraditionsPref = source["TraditionsPref"];
        this.FacilitiesPref = source["FacilitiesPref"];
        this.AtmospherePref = source["AtmospherePref"];
        this.AcademicsPref = source["AcademicsPref"];
        this.ConferencePref = source["ConferencePref"];
        this.CoachPref = source["CoachPref"];
        this.SeasonMomentumPref = source["SeasonMomentumPref"];
        this.LeadingTeams = this.convertValues(source["LeadingTeams"], LeadingTeams);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CollegeStandings {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    TeamName: string;
    SeasonID: number;
    Season: number;
    LeagueID: number;
    ConferenceID: number;
    TotalWins: number;
    TotalLosses: number;
    TotalOTWins: number;
    TotalOTLosses: number;
    ShootoutWins: number;
    ShootoutLosses: number;
    ConferenceWins: number;
    ConferenceLosses: number;
    ConferenceOTWins: number;
    ConferenceOTLosses: number;
    RankedWins: number;
    RankedLosses: number;
    Points: number;
    GoalsFor: number;
    GoalsAgainst: number;
    Streak: number;
    IsWinStreak: boolean;
    HomeWins: number;
    AwayWins: number;
    PostSeasonStatus: string;
    Rank: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.TeamName = source["TeamName"];
        this.SeasonID = source["SeasonID"];
        this.Season = source["Season"];
        this.LeagueID = source["LeagueID"];
        this.ConferenceID = source["ConferenceID"];
        this.TotalWins = source["TotalWins"];
        this.TotalLosses = source["TotalLosses"];
        this.TotalOTWins = source["TotalOTWins"];
        this.TotalOTLosses = source["TotalOTLosses"];
        this.ShootoutWins = source["ShootoutWins"];
        this.ShootoutLosses = source["ShootoutLosses"];
        this.ConferenceWins = source["ConferenceWins"];
        this.ConferenceLosses = source["ConferenceLosses"];
        this.ConferenceOTWins = source["ConferenceOTWins"];
        this.ConferenceOTLosses = source["ConferenceOTLosses"];
        this.RankedWins = source["RankedWins"];
        this.RankedLosses = source["RankedLosses"];
        this.Points = source["Points"];
        this.GoalsFor = source["GoalsFor"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.Streak = source["Streak"];
        this.IsWinStreak = source["IsWinStreak"];
        this.HomeWins = source["HomeWins"];
        this.AwayWins = source["AwayWins"];
        this.PostSeasonStatus = source["PostSeasonStatus"];
        this.Rank = source["Rank"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class BootstrapData {
    CollegeTeam: CollegeTeam;
    AllCollegeTeams: CollegeTeam[];
    CollegeStandings: CollegeStandings[];
    CollegeRosterMap: {[key: uint]: CollegePlayer[]};
    Recruits: Croot[];
    RecruitProfiles: RecruitPlayerProfile[];
    TeamProfileMap: {[key: uint]: RecruitingTeamProfile};
    PortalPlayers: CollegePlayer[];
    CollegeInjuryReport: CollegePlayer[];
    CollegeNews: NewsLog[];
    CollegeNotifications: Notification[];
    AllCollegeGames: CollegeGame[];
    CHLGameplan: CollegeGameplan;
    CollegeTeamLineups: CollegeLineup[];
    CollegeTeamShootoutLineup: CollegeShootoutLineup;
    TopCHLGoals: CollegePlayer[];
    TopCHLAssists: CollegePlayer[];
    TopCHLSaves: CollegePlayer[];
    CollegePoll: CollegePollSubmission;
    OfficialPolls: CollegePollOfficial[];
    ProTeam: ProfessionalTeam;
    AllProTeams: ProfessionalTeam[];
    ProStandings: ProfessionalStandings[];
    ProRosterMap: {[key: uint]: ProfessionalPlayer[]};
    AffiliatePlayers: ProfessionalPlayer[];
    TopPHLGoals: ProfessionalPlayer[];
    TopPHLAssists: ProfessionalPlayer[];
    TopPHLSaves: ProfessionalPlayer[];
    FreeAgentOffers: FreeAgencyOffer[];
    WaiverWireOffers: WaiverOffer[];
    CapsheetMap: {[key: uint]: ProCapsheet};
    ProInjuryReport: ProfessionalPlayer[];
    ProNews: NewsLog[];
    ProNotifications: Notification[];
    AllProGames: ProfessionalGame[];
    PHLGameplan: ProGameplan;
    ProTeamLineups: ProfessionalLineup[];
    ProTeamShootoutLineup: ProfessionalShootoutLineup;
    FaceData: {[key: uint]: FaceDataResponse};
    ContractMap: {[key: uint]: ProContract};
    ExtensionMap: {[key: uint]: ExtensionOffer};
    ProTradeProposalMap: {[key: uint]: TradeProposal[]};
    ProTradePreferenceMap: {[key: uint]: TradePreferences};
    DraftPicks: DraftPick[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.CollegeTeam = this.convertValues(source["CollegeTeam"], CollegeTeam);
        this.AllCollegeTeams = this.convertValues(source["AllCollegeTeams"], CollegeTeam);
        this.CollegeStandings = this.convertValues(source["CollegeStandings"], CollegeStandings);
        this.CollegeRosterMap = source["CollegeRosterMap"];
        this.Recruits = this.convertValues(source["Recruits"], Croot);
        this.RecruitProfiles = this.convertValues(source["RecruitProfiles"], RecruitPlayerProfile);
        this.TeamProfileMap = source["TeamProfileMap"];
        this.PortalPlayers = this.convertValues(source["PortalPlayers"], CollegePlayer);
        this.CollegeInjuryReport = this.convertValues(source["CollegeInjuryReport"], CollegePlayer);
        this.CollegeNews = this.convertValues(source["CollegeNews"], NewsLog);
        this.CollegeNotifications = this.convertValues(source["CollegeNotifications"], Notification);
        this.AllCollegeGames = this.convertValues(source["AllCollegeGames"], CollegeGame);
        this.CHLGameplan = this.convertValues(source["CHLGameplan"], CollegeGameplan);
        this.CollegeTeamLineups = this.convertValues(source["CollegeTeamLineups"], CollegeLineup);
        this.CollegeTeamShootoutLineup = this.convertValues(source["CollegeTeamShootoutLineup"], CollegeShootoutLineup);
        this.TopCHLGoals = this.convertValues(source["TopCHLGoals"], CollegePlayer);
        this.TopCHLAssists = this.convertValues(source["TopCHLAssists"], CollegePlayer);
        this.TopCHLSaves = this.convertValues(source["TopCHLSaves"], CollegePlayer);
        this.CollegePoll = this.convertValues(source["CollegePoll"], CollegePollSubmission);
        this.OfficialPolls = this.convertValues(source["OfficialPolls"], CollegePollOfficial);
        this.ProTeam = this.convertValues(source["ProTeam"], ProfessionalTeam);
        this.AllProTeams = this.convertValues(source["AllProTeams"], ProfessionalTeam);
        this.ProStandings = this.convertValues(source["ProStandings"], ProfessionalStandings);
        this.ProRosterMap = source["ProRosterMap"];
        this.AffiliatePlayers = this.convertValues(source["AffiliatePlayers"], ProfessionalPlayer);
        this.TopPHLGoals = this.convertValues(source["TopPHLGoals"], ProfessionalPlayer);
        this.TopPHLAssists = this.convertValues(source["TopPHLAssists"], ProfessionalPlayer);
        this.TopPHLSaves = this.convertValues(source["TopPHLSaves"], ProfessionalPlayer);
        this.FreeAgentOffers = this.convertValues(source["FreeAgentOffers"], FreeAgencyOffer);
        this.WaiverWireOffers = this.convertValues(source["WaiverWireOffers"], WaiverOffer);
        this.CapsheetMap = this.convertValues(source["CapsheetMap"], ProCapsheet, true);
        this.ProInjuryReport = this.convertValues(source["ProInjuryReport"], ProfessionalPlayer);
        this.ProNews = this.convertValues(source["ProNews"], NewsLog);
        this.ProNotifications = this.convertValues(source["ProNotifications"], Notification);
        this.AllProGames = this.convertValues(source["AllProGames"], ProfessionalGame);
        this.PHLGameplan = this.convertValues(source["PHLGameplan"], ProGameplan);
        this.ProTeamLineups = this.convertValues(source["ProTeamLineups"], ProfessionalLineup);
        this.ProTeamShootoutLineup = this.convertValues(source["ProTeamShootoutLineup"], ProfessionalShootoutLineup);
        this.FaceData = this.convertValues(source["FaceData"], FaceDataResponse, true);
        this.ContractMap = this.convertValues(source["ContractMap"], ProContract, true);
        this.ExtensionMap = this.convertValues(source["ExtensionMap"], ExtensionOffer, true);
        this.ProTradeProposalMap = source["ProTradeProposalMap"];
        this.ProTradePreferenceMap = this.convertValues(source["ProTradePreferenceMap"], TradePreferences, true);
        this.DraftPicks = this.convertValues(source["DraftPicks"], DraftPick);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}





export class BaseRecruitingGrades {
    AgilityGrade: string;
    FaceoffsGrade: string;
    CloseShotAccuracyGrade: string;
    CloseShotPowerGrade: string;
    LongShotAccuracyGrade: string;
    LongShotPowerGrade: string;
    PassingGrade: string;
    PuckHandlingGrade: string;
    StrengthGrade: string;
    BodyCheckingGrade: string;
    StickCheckingGrade: string;
    ShotBlockingGrade: string;
    GoalkeepingGrade: string;
    GoalieVisionGrade: string;
    GoalieReboundGrade: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.AgilityGrade = source["AgilityGrade"];
        this.FaceoffsGrade = source["FaceoffsGrade"];
        this.CloseShotAccuracyGrade = source["CloseShotAccuracyGrade"];
        this.CloseShotPowerGrade = source["CloseShotPowerGrade"];
        this.LongShotAccuracyGrade = source["LongShotAccuracyGrade"];
        this.LongShotPowerGrade = source["LongShotPowerGrade"];
        this.PassingGrade = source["PassingGrade"];
        this.PuckHandlingGrade = source["PuckHandlingGrade"];
        this.StrengthGrade = source["StrengthGrade"];
        this.BodyCheckingGrade = source["BodyCheckingGrade"];
        this.StickCheckingGrade = source["StickCheckingGrade"];
        this.ShotBlockingGrade = source["ShotBlockingGrade"];
        this.GoalkeepingGrade = source["GoalkeepingGrade"];
        this.GoalieVisionGrade = source["GoalieVisionGrade"];
        this.GoalieReboundGrade = source["GoalieReboundGrade"];
    }
}

export class CreateRecruitProfileDto {
    PlayerID: number;
    SeasonID: number;
    RecruitID: number;
    ProfileID: number;
    Team: CollegeTeam;
    PlayerRecruit: Croot;
    Recruiter: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.PlayerID = source["PlayerID"];
        this.SeasonID = source["SeasonID"];
        this.RecruitID = source["RecruitID"];
        this.ProfileID = source["ProfileID"];
        this.Team = this.convertValues(source["Team"], CollegeTeam);
        this.PlayerRecruit = this.convertValues(source["PlayerRecruit"], Croot);
        this.Recruiter = source["Recruiter"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class UpdateRecruitProfileDto {
    RecruitPointsID: number;
    RecruitID: number;
    ProfileID: number;
    Team: string;
    WeekID: number;
    AllocationID: number;
    SpentPoints: number;
    RewardScholarship: boolean;
    RevokeScholarship: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.RecruitPointsID = source["RecruitPointsID"];
        this.RecruitID = source["RecruitID"];
        this.ProfileID = source["ProfileID"];
        this.Team = source["Team"];
        this.WeekID = source["WeekID"];
        this.AllocationID = source["AllocationID"];
        this.SpentPoints = source["SpentPoints"];
        this.RewardScholarship = source["RewardScholarship"];
        this.RevokeScholarship = source["RevokeScholarship"];
    }
}
export class ScoutAttributeDTO {
    ProfileID: number;
    RecruitID: number;
    Attribute: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ProfileID = source["ProfileID"];
        this.RecruitID = source["RecruitID"];
        this.Attribute = source["Attribute"];
    }
}
export class CrootProfile {
    ID: number;
    SeasonID: number;
    RecruitID: number;
    ProfileID: number;
    TotalPoints: number;
    CurrentWeeksPoints: number;
    Modifier: number;
    SpendingCount: number;
    Scholarship: boolean;
    ScholarshipRevoked: boolean;
    IsCloseToHome: boolean;
    IsPipeline: boolean;
    TeamAbbreviation: string;
    RemovedFromBoard: boolean;
    IsSigned: boolean;
    IsLocked: boolean;
    CaughtCheating: boolean;
    Agility: boolean;
    Faceoffs: boolean;
    LongShotAccuracy: boolean;
    LongShotPower: boolean;
    CloseShotAccuracy: boolean;
    CloseShotPower: boolean;
    OneTimer: boolean;
    Passing: boolean;
    PuckHandling: boolean;
    Strength: boolean;
    BodyChecking: boolean;
    StickChecking: boolean;
    ShotBlocking: boolean;
    Goalkeeping: boolean;
    GoalieVision: boolean;
    GoalieReboundControl: boolean;
    Recruit: Croot;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.SeasonID = source["SeasonID"];
        this.RecruitID = source["RecruitID"];
        this.ProfileID = source["ProfileID"];
        this.TotalPoints = source["TotalPoints"];
        this.CurrentWeeksPoints = source["CurrentWeeksPoints"];
        this.Modifier = source["Modifier"];
        this.SpendingCount = source["SpendingCount"];
        this.Scholarship = source["Scholarship"];
        this.ScholarshipRevoked = source["ScholarshipRevoked"];
        this.IsCloseToHome = source["IsCloseToHome"];
        this.IsPipeline = source["IsPipeline"];
        this.TeamAbbreviation = source["TeamAbbreviation"];
        this.RemovedFromBoard = source["RemovedFromBoard"];
        this.IsSigned = source["IsSigned"];
        this.IsLocked = source["IsLocked"];
        this.CaughtCheating = source["CaughtCheating"];
        this.Agility = source["Agility"];
        this.Faceoffs = source["Faceoffs"];
        this.LongShotAccuracy = source["LongShotAccuracy"];
        this.LongShotPower = source["LongShotPower"];
        this.CloseShotAccuracy = source["CloseShotAccuracy"];
        this.CloseShotPower = source["CloseShotPower"];
        this.OneTimer = source["OneTimer"];
        this.Passing = source["Passing"];
        this.PuckHandling = source["PuckHandling"];
        this.Strength = source["Strength"];
        this.BodyChecking = source["BodyChecking"];
        this.StickChecking = source["StickChecking"];
        this.ShotBlocking = source["ShotBlocking"];
        this.Goalkeeping = source["Goalkeeping"];
        this.GoalieVision = source["GoalieVision"];
        this.GoalieReboundControl = source["GoalieReboundControl"];
        this.Recruit = this.convertValues(source["Recruit"], Croot);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class SimTeamBoardResponse {
    ID: number;
    TeamID: number;
    Team: string;
    State: string;
    ScholarshipsAvailable: number;
    WeeklyPoints: number;
    SpentPoints: number;
    TotalCommitments: number;
    RecruitClassSize: number;
    BaseEfficiencyScore: number;
    RecruitingEfficiencyScore: number;
    PreviousOverallWinPer: number;
    PreviousConferenceWinPer: number;
    CurrentOverallWinPer: number;
    CurrentConferenceWinPer: number;
    ESPNScore: number;
    RivalsScore: number;
    Rank247Score: number;
    CompositeScore: number;
    IsAI: boolean;
    IsUserTeam: boolean;
    BattlesWon: number;
    BattlesLost: number;
    AIMinThreshold: number;
    AIMaxThreshold: number;
    AIStarMin: number;
    AIStarMax: number;
    AIAutoOfferscholarships: boolean;
    OffensiveScheme: string;
    DefensiveScheme: string;
    Recruiter: string;
    RecruitingClassRank: number;
    Recruits: CrootProfile[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.State = source["State"];
        this.ScholarshipsAvailable = source["ScholarshipsAvailable"];
        this.WeeklyPoints = source["WeeklyPoints"];
        this.SpentPoints = source["SpentPoints"];
        this.TotalCommitments = source["TotalCommitments"];
        this.RecruitClassSize = source["RecruitClassSize"];
        this.BaseEfficiencyScore = source["BaseEfficiencyScore"];
        this.RecruitingEfficiencyScore = source["RecruitingEfficiencyScore"];
        this.PreviousOverallWinPer = source["PreviousOverallWinPer"];
        this.PreviousConferenceWinPer = source["PreviousConferenceWinPer"];
        this.CurrentOverallWinPer = source["CurrentOverallWinPer"];
        this.CurrentConferenceWinPer = source["CurrentConferenceWinPer"];
        this.ESPNScore = source["ESPNScore"];
        this.RivalsScore = source["RivalsScore"];
        this.Rank247Score = source["Rank247Score"];
        this.CompositeScore = source["CompositeScore"];
        this.IsAI = source["IsAI"];
        this.IsUserTeam = source["IsUserTeam"];
        this.BattlesWon = source["BattlesWon"];
        this.BattlesLost = source["BattlesLost"];
        this.AIMinThreshold = source["AIMinThreshold"];
        this.AIMaxThreshold = source["AIMaxThreshold"];
        this.AIStarMin = source["AIStarMin"];
        this.AIStarMax = source["AIStarMax"];
        this.AIAutoOfferscholarships = source["AIAutoOfferscholarships"];
        this.OffensiveScheme = source["OffensiveScheme"];
        this.DefensiveScheme = source["DefensiveScheme"];
        this.Recruiter = source["Recruiter"];
        this.RecruitingClassRank = source["RecruitingClassRank"];
        this.Recruits = this.convertValues(source["Recruits"], CrootProfile);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class UpdateRecruitingBoardDTO {
    Profile: RecruitingTeamProfile;
    Recruits: RecruitPlayerProfile[];
    TeamID: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Profile = this.convertValues(source["Profile"], RecruitingTeamProfile);
        this.Recruits = this.convertValues(source["Recruits"], RecruitPlayerProfile);
        this.TeamID = source["TeamID"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class RecruitPointAllocation {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    RecruitID: number;
    TeamProfileID: number;
    RecruitProfileID: number;
    WeekID: number;
    Points: number;
    ModAffectedPoints: number;
    IsHomeStateApplied: boolean;
    IsPipelineApplied: boolean;
    CaughtCheating: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.RecruitID = source["RecruitID"];
        this.TeamProfileID = source["TeamProfileID"];
        this.RecruitProfileID = source["RecruitProfileID"];
        this.WeekID = source["WeekID"];
        this.Points = source["Points"];
        this.ModAffectedPoints = source["ModAffectedPoints"];
        this.IsHomeStateApplied = source["IsHomeStateApplied"];
        this.IsPipelineApplied = source["IsPipelineApplied"];
        this.CaughtCheating = source["CaughtCheating"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class RecruitingOdds {
    Odds: number;
    IsCloseToHome: boolean;
    IsPipeline: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Odds = source["Odds"];
        this.IsCloseToHome = source["IsCloseToHome"];
        this.IsPipeline = source["IsPipeline"];
    }
}












export class TeamRequest {
    TeamID: number;
    Username: string;
    IsApproved: boolean;
    IsActive: boolean;
    Role: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.TeamID = source["TeamID"];
        this.Username = source["Username"];
        this.IsApproved = source["IsApproved"];
        this.IsActive = source["IsActive"];
        this.Role = source["Role"];
    }
}
export class CollegeTeamRequest {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    Username: string;
    IsApproved: boolean;
    IsActive: boolean;
    Role: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.Username = source["Username"];
        this.IsApproved = source["IsApproved"];
        this.IsActive = source["IsActive"];
        this.Role = source["Role"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class ProTeamRequest {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    TeamID: number;
    Username: string;
    IsApproved: boolean;
    IsActive: boolean;
    Role: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.TeamID = source["TeamID"];
        this.Username = source["Username"];
        this.IsApproved = source["IsApproved"];
        this.IsActive = source["IsActive"];
        this.Role = source["Role"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}


export class PollDataResponse {
    Poll: CollegePollSubmission;
    Matches: CollegeGame[];
    OfficialPolls: CollegePollOfficial[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.Poll = this.convertValues(source["Poll"], CollegePollSubmission);
        this.Matches = this.convertValues(source["Matches"], CollegeGame);
        this.OfficialPolls = this.convertValues(source["OfficialPolls"], CollegePollOfficial);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}

export class FreeAgencyOfferDTO {
    ID: number;
    PlayerID: number;
    TeamID: number;
    Team: string;
    ContractLength: number;
    Y1BaseSalary: number;
    Y2BaseSalary: number;
    Y3BaseSalary: number;
    Y4BaseSalary: number;
    Y5BaseSalary: number;
    ToAffiliate: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.ContractLength = source["ContractLength"];
        this.Y1BaseSalary = source["Y1BaseSalary"];
        this.Y2BaseSalary = source["Y2BaseSalary"];
        this.Y3BaseSalary = source["Y3BaseSalary"];
        this.Y4BaseSalary = source["Y4BaseSalary"];
        this.Y5BaseSalary = source["Y5BaseSalary"];
        this.ToAffiliate = source["ToAffiliate"];
    }
}

export class WaiverOfferDTO {
    ID: number;
    PlayerID: number;
    TeamID: number;
    Team: string;
    WaiverOrder: number;
    IsActive: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.WaiverOrder = source["WaiverOrder"];
        this.IsActive = source["IsActive"];
    }
}
export class ShootoutPlayerIDs {
    ID: number;
    TeamID: number;
    Shooter1ID: number;
    Shooter1ShotType: number;
    Shooter2ID: number;
    Shooter2ShotType: number;
    Shooter3ID: number;
    Shooter3ShotType: number;
    Shooter4ID: number;
    Shooter4ShotType: number;
    Shooter5ID: number;
    Shooter5ShotType: number;
    Shooter6ID: number;
    Shooter6ShotType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.TeamID = source["TeamID"];
        this.Shooter1ID = source["Shooter1ID"];
        this.Shooter1ShotType = source["Shooter1ShotType"];
        this.Shooter2ID = source["Shooter2ID"];
        this.Shooter2ShotType = source["Shooter2ShotType"];
        this.Shooter3ID = source["Shooter3ID"];
        this.Shooter3ShotType = source["Shooter3ShotType"];
        this.Shooter4ID = source["Shooter4ID"];
        this.Shooter4ShotType = source["Shooter4ShotType"];
        this.Shooter5ID = source["Shooter5ID"];
        this.Shooter5ShotType = source["Shooter5ShotType"];
        this.Shooter6ID = source["Shooter6ID"];
        this.Shooter6ShotType = source["Shooter6ShotType"];
    }
}
export class Allocations {
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
    }
}
export class BaseLineup {
    TeamID: number;
    Line: number;
    LineType: number;
    CenterID: number;
    Forward1ID: number;
    Forward2ID: number;
    Defender1ID: number;
    Defender2ID: number;
    GoalieID: number;
    AGZShot: number;
    AGZPass: number;
    AGZPassBack: number;
    AGZAgility: number;
    AGZStickCheck: number;
    AGZBodyCheck: number;
    AZShot: number;
    AZPass: number;
    AZLongPass: number;
    AZAgility: number;
    AZStickCheck: number;
    AZBodyCheck: number;
    NPass: number;
    NAgility: number;
    NStickCheck: number;
    NBodyCheck: number;
    DZPass: number;
    DZPassBack: number;
    DZAgility: number;
    DZStickCheck: number;
    DZBodyCheck: number;
    DGZPass: number;
    DGZLongPass: number;
    DGZAgility: number;
    DGZStickCheck: number;
    DGZBodyCheck: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.TeamID = source["TeamID"];
        this.Line = source["Line"];
        this.LineType = source["LineType"];
        this.CenterID = source["CenterID"];
        this.Forward1ID = source["Forward1ID"];
        this.Forward2ID = source["Forward2ID"];
        this.Defender1ID = source["Defender1ID"];
        this.Defender2ID = source["Defender2ID"];
        this.GoalieID = source["GoalieID"];
        this.AGZShot = source["AGZShot"];
        this.AGZPass = source["AGZPass"];
        this.AGZPassBack = source["AGZPassBack"];
        this.AGZAgility = source["AGZAgility"];
        this.AGZStickCheck = source["AGZStickCheck"];
        this.AGZBodyCheck = source["AGZBodyCheck"];
        this.AZShot = source["AZShot"];
        this.AZPass = source["AZPass"];
        this.AZLongPass = source["AZLongPass"];
        this.AZAgility = source["AZAgility"];
        this.AZStickCheck = source["AZStickCheck"];
        this.AZBodyCheck = source["AZBodyCheck"];
        this.NPass = source["NPass"];
        this.NAgility = source["NAgility"];
        this.NStickCheck = source["NStickCheck"];
        this.NBodyCheck = source["NBodyCheck"];
        this.DZPass = source["DZPass"];
        this.DZPassBack = source["DZPassBack"];
        this.DZAgility = source["DZAgility"];
        this.DZStickCheck = source["DZStickCheck"];
        this.DZBodyCheck = source["DZBodyCheck"];
        this.DGZPass = source["DGZPass"];
        this.DGZLongPass = source["DGZLongPass"];
        this.DGZAgility = source["DGZAgility"];
        this.DGZStickCheck = source["DGZStickCheck"];
        this.DGZBodyCheck = source["DGZBodyCheck"];
    }
}
export class TopPlayer {
    PlayerID: number;
    FirstName: string;
    LastName: string;
    Position: string;
    Archetype: string;
    OverallGrade: string;
    Overall: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.PlayerID = source["PlayerID"];
        this.FirstName = source["FirstName"];
        this.LastName = source["LastName"];
        this.Position = source["Position"];
        this.Archetype = source["Archetype"];
        this.OverallGrade = source["OverallGrade"];
        this.Overall = source["Overall"];
    }
}
export class TeamRecordResponse {
    OverallWins: number;
    OverallLosses: number;
    CurrentSeasonWins: number;
    CurrentSeasonLosses: number;
    PostSeasonWins: number;
    PostSeasonLosses: number;
    ConferenceChampionships: string[];
    DivisionTitles: string[];
    NationalChampionships: string[];
    TopPlayers: TopPlayer[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.OverallWins = source["OverallWins"];
        this.OverallLosses = source["OverallLosses"];
        this.CurrentSeasonWins = source["CurrentSeasonWins"];
        this.CurrentSeasonLosses = source["CurrentSeasonLosses"];
        this.PostSeasonWins = source["PostSeasonWins"];
        this.PostSeasonLosses = source["PostSeasonLosses"];
        this.ConferenceChampionships = source["ConferenceChampionships"];
        this.DivisionTitles = source["DivisionTitles"];
        this.NationalChampionships = source["NationalChampionships"];
        this.TopPlayers = this.convertValues(source["TopPlayers"], TopPlayer);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}

export class InboxResponse {
    CHLNotifications: Notification[];
    PHLNotifications: Notification[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.CHLNotifications = this.convertValues(source["CHLNotifications"], Notification);
        this.PHLNotifications = this.convertValues(source["PHLNotifications"], Notification);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class BasePlayerStats {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    StartedGame: boolean;
    GameDay: string;
    PlayerID: number;
    TeamID: number;
    SeasonID: number;
    Team: string;
    PreviousTeamID: number;
    PreviousTeam: string;
    Goals: number;
    Assists: number;
    Points: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    TimeOnIce: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    GoalieWins: number;
    GoalieLosses: number;
    GoalieTies: number;
    OvertimeLosses: number;
    ShotsAgainst: number;
    Saves: number;
    GoalsAgainst: number;
    SavePercentage: number;
    Shutouts: number;
    ShotsBlocked: number;
    BodyChecks: number;
    StickChecks: number;
    IsInjured: boolean;
    DaysOfRecovery: number;
    InjuryName: string;
    InjuryType: string;
    GameType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.StartedGame = source["StartedGame"];
        this.GameDay = source["GameDay"];
        this.PlayerID = source["PlayerID"];
        this.TeamID = source["TeamID"];
        this.SeasonID = source["SeasonID"];
        this.Team = source["Team"];
        this.PreviousTeamID = source["PreviousTeamID"];
        this.PreviousTeam = source["PreviousTeam"];
        this.Goals = source["Goals"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.TimeOnIce = source["TimeOnIce"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.GoalieWins = source["GoalieWins"];
        this.GoalieLosses = source["GoalieLosses"];
        this.GoalieTies = source["GoalieTies"];
        this.OvertimeLosses = source["OvertimeLosses"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.ShotsBlocked = source["ShotsBlocked"];
        this.BodyChecks = source["BodyChecks"];
        this.StickChecks = source["StickChecks"];
        this.IsInjured = source["IsInjured"];
        this.DaysOfRecovery = source["DaysOfRecovery"];
        this.InjuryName = source["InjuryName"];
        this.InjuryType = source["InjuryType"];
        this.GameType = source["GameType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class BaseTeamStats {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    SeasonID: number;
    TeamID: number;
    Team: string;
    GameDay: string;
    GoalsFor: number;
    GoalsAgainst: number;
    Assists: number;
    Points: number;
    Period1Score: number;
    Period2Score: number;
    Period3Score: number;
    OTScore: number;
    PlusMinus: number;
    PenaltyMinutes: number;
    EvenStrengthGoals: number;
    EvenStrengthPoints: number;
    PowerPlayGoals: number;
    PowerPlayPoints: number;
    ShorthandedGoals: number;
    ShorthandedPoints: number;
    OvertimeGoals: number;
    GameWinningGoals: number;
    Shots: number;
    ShootingPercentage: number;
    FaceOffWinPercentage: number;
    FaceOffsWon: number;
    FaceOffs: number;
    ShotsAgainst: number;
    Saves: number;
    SavePercentage: number;
    Shutouts: number;
    GameType: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.SeasonID = source["SeasonID"];
        this.TeamID = source["TeamID"];
        this.Team = source["Team"];
        this.GameDay = source["GameDay"];
        this.GoalsFor = source["GoalsFor"];
        this.GoalsAgainst = source["GoalsAgainst"];
        this.Assists = source["Assists"];
        this.Points = source["Points"];
        this.Period1Score = source["Period1Score"];
        this.Period2Score = source["Period2Score"];
        this.Period3Score = source["Period3Score"];
        this.OTScore = source["OTScore"];
        this.PlusMinus = source["PlusMinus"];
        this.PenaltyMinutes = source["PenaltyMinutes"];
        this.EvenStrengthGoals = source["EvenStrengthGoals"];
        this.EvenStrengthPoints = source["EvenStrengthPoints"];
        this.PowerPlayGoals = source["PowerPlayGoals"];
        this.PowerPlayPoints = source["PowerPlayPoints"];
        this.ShorthandedGoals = source["ShorthandedGoals"];
        this.ShorthandedPoints = source["ShorthandedPoints"];
        this.OvertimeGoals = source["OvertimeGoals"];
        this.GameWinningGoals = source["GameWinningGoals"];
        this.Shots = source["Shots"];
        this.ShootingPercentage = source["ShootingPercentage"];
        this.FaceOffWinPercentage = source["FaceOffWinPercentage"];
        this.FaceOffsWon = source["FaceOffsWon"];
        this.FaceOffs = source["FaceOffs"];
        this.ShotsAgainst = source["ShotsAgainst"];
        this.Saves = source["Saves"];
        this.SavePercentage = source["SavePercentage"];
        this.Shutouts = source["Shutouts"];
        this.GameType = source["GameType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class TeamSeasonStats {
    StatType: number;
    GamesPlayed: number;
    PointPercentage: number;
    RegulationWins: number;
    RegulationOvertimeWins: number;
    ShootoutsWon: number;
    GFGP: number;
    GAGP: number;
    PowerPlayPercentage: number;
    PenaltyKillPercentage: number;
    PowerPlayNetPercentage: number;
    PenaltyKillNetPercentage: number;
    ShotsPerGame: number;
    ShotsAgainstPerGame: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.StatType = source["StatType"];
        this.GamesPlayed = source["GamesPlayed"];
        this.PointPercentage = source["PointPercentage"];
        this.RegulationWins = source["RegulationWins"];
        this.RegulationOvertimeWins = source["RegulationOvertimeWins"];
        this.ShootoutsWon = source["ShootoutsWon"];
        this.GFGP = source["GFGP"];
        this.GAGP = source["GAGP"];
        this.PowerPlayPercentage = source["PowerPlayPercentage"];
        this.PenaltyKillPercentage = source["PenaltyKillPercentage"];
        this.PowerPlayNetPercentage = source["PowerPlayNetPercentage"];
        this.PenaltyKillNetPercentage = source["PenaltyKillNetPercentage"];
        this.ShotsPerGame = source["ShotsPerGame"];
        this.ShotsAgainstPerGame = source["ShotsAgainstPerGame"];
    }
}










export class Timestamp {
    ID: number;
    CreatedAt: Time;
    UpdatedAt: Time;
    DeletedAt: DeletedAt;
    RunCron: boolean;
    RunGames: boolean;
    WeekID: number;
    Week: number;
    SeasonID: number;
    Season: number;
    GamesARan: boolean;
    GamesBRan: boolean;
    GamesCRan: boolean;
    GamesDRan: boolean;
    CollegePollRan: boolean;
    RecruitingSynced: boolean;
    GMActionsCompleted: boolean;
    IsOffSeason: boolean;
    IsRecruitingLocked: boolean;
    AIDepthchartsSync: boolean;
    AIRecruitingBoardsSynced: boolean;
    CollegeSeasonOver: boolean;
    NHLSeasonOver: boolean;
    CrootsGenerated: boolean;
    ProgressedCollegePlayers: boolean;
    ProgressedProfessionalPlayers: boolean;
    TransferPortalPhase: number;
    TransferPortalRound: number;
    IsFreeAgencyLocked: boolean;
    FreeAgencyRound: number;
    IsDraftTime: boolean;
    Y1Capspace: number;
    Y2Capspace: number;
    Y3Capspace: number;
    Y4Capspace: number;
    Y5Capspace: number;
    DeadCapLimit: number;
    PreseasonPhase: number;
    IsPreseason: boolean;
    SeasonPhase: number;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.ID = source["ID"];
        this.CreatedAt = this.convertValues(source["CreatedAt"], Time);
        this.UpdatedAt = this.convertValues(source["UpdatedAt"], Time);
        this.DeletedAt = this.convertValues(source["DeletedAt"], DeletedAt);
        this.RunCron = source["RunCron"];
        this.RunGames = source["RunGames"];
        this.WeekID = source["WeekID"];
        this.Week = source["Week"];
        this.SeasonID = source["SeasonID"];
        this.Season = source["Season"];
        this.GamesARan = source["GamesARan"];
        this.GamesBRan = source["GamesBRan"];
        this.GamesCRan = source["GamesCRan"];
        this.GamesDRan = source["GamesDRan"];
        this.CollegePollRan = source["CollegePollRan"];
        this.RecruitingSynced = source["RecruitingSynced"];
        this.GMActionsCompleted = source["GMActionsCompleted"];
        this.IsOffSeason = source["IsOffSeason"];
        this.IsRecruitingLocked = source["IsRecruitingLocked"];
        this.AIDepthchartsSync = source["AIDepthchartsSync"];
        this.AIRecruitingBoardsSynced = source["AIRecruitingBoardsSynced"];
        this.CollegeSeasonOver = source["CollegeSeasonOver"];
        this.NHLSeasonOver = source["NHLSeasonOver"];
        this.CrootsGenerated = source["CrootsGenerated"];
        this.ProgressedCollegePlayers = source["ProgressedCollegePlayers"];
        this.ProgressedProfessionalPlayers = source["ProgressedProfessionalPlayers"];
        this.TransferPortalPhase = source["TransferPortalPhase"];
        this.TransferPortalRound = source["TransferPortalRound"];
        this.IsFreeAgencyLocked = source["IsFreeAgencyLocked"];
        this.FreeAgencyRound = source["FreeAgencyRound"];
        this.IsDraftTime = source["IsDraftTime"];
        this.Y1Capspace = source["Y1Capspace"];
        this.Y2Capspace = source["Y2Capspace"];
        this.Y3Capspace = source["Y3Capspace"];
        this.Y4Capspace = source["Y4Capspace"];
        this.Y5Capspace = source["Y5Capspace"];
        this.DeadCapLimit = source["DeadCapLimit"];
        this.PreseasonPhase = source["PreseasonPhase"];
        this.IsPreseason = source["IsPreseason"];
        this.SeasonPhase = source["SeasonPhase"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class TeamRequestsResponse {
    CollegeRequests: CollegeTeamRequest[];
    ProRequest: ProTeamRequest[];
    AcceptedTrades: TradeProposal[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.CollegeRequests = this.convertValues(source["CollegeRequests"], CollegeTeamRequest);
        this.ProRequest = this.convertValues(source["ProRequest"], ProTeamRequest);
        this.AcceptedTrades = this.convertValues(source["AcceptedTrades"], TradeProposal);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (Array.isArray(a)) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}