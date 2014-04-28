package core;

public interface GameChangeListener {
	public void selectQuestion(Question q,int row,int col);
	public void answeredQuestion(int playerID,int score, int choice) throws InterruptedException;
	public void buzzedIn(int playerID);

}
