package core;

public interface GameChangeListener {
	public void selectQuestion(Question q,int row,int col);
	public void answeredQuestion(int playerID,int score);
	public void buzzedIn(int playerID);

}
