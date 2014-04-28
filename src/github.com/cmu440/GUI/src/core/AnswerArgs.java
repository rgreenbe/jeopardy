package core;

public class AnswerArgs {
	private int playerID;
	private int gameID;
	private int scoreChange;
	private int choice;
	public AnswerArgs(int gameID,int playerID,int scoreChange,int choice){
		this.playerID=playerID;
		this.gameID=gameID;
		this.scoreChange=scoreChange;
		this.choice=choice;
		
	}
	public int scoreChange() {
		return scoreChange;
	}
	public int playerID() {
		return playerID;
	}
	public int choice(){
		return choice;
	}
}
