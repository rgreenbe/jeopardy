package gui;

import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.IOException;

import javax.swing.JButton;

import core.Jeopardy;

public class JoinListener implements ActionListener{
	Jeopardy game;
	JButton join;
	public JoinListener(Jeopardy game, JButton join){
		this.game=game;
		this.join=join;
	}

	@Override
	public void actionPerformed(ActionEvent arg0) {
		try {
			join.setEnabled(false);
			game.joinGame();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}		
	}

}
