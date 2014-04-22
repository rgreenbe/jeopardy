package gui;

import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;

import core.Jeopardy;

public class JoinListener implements ActionListener{
	Jeopardy game;
	public JoinListener(Jeopardy game){
		this.game=game;
	}

	@Override
	public void actionPerformed(ActionEvent arg0) {
		game.joinGame();		
	}

}
