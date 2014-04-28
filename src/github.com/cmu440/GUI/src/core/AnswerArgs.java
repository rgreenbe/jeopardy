package core;

public class AnswerArgs {
	private int playerID;
	private int gameID;
	private int scoreChange;
	public AnswerArgs(int gameID,int playerID,int scoreChange){
		this.playerID=playerID;
		this.gameID=gameID;
		this.scoreChange=scoreChange;
		
	}
	public int scoreChange() {
		return scoreChange;
	}
	public int playerID() {
		return playerID;
	}
}
