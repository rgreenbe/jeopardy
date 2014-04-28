package core;

import java.io.FileNotFoundException;

import gui.JeopardyPanel;

import javax.swing.JFrame;
import javax.swing.SwingUtilities;

public class Main {
	  private static final String NAME = "Scrabble";

	    public static void main(String[] args) {
	        SwingUtilities.invokeLater(new Runnable() {
	            @Override
	            public void run() {
	                try {
						createAndShowJeopardy();
					} catch (FileNotFoundException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					}
	            }
	        });
	    }
	    private static void createAndShowJeopardy() throws FileNotFoundException {
	        // Create and set-up the window.
	        JFrame frame = new JFrame(NAME);
	        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);

	        Jeopardy game = new Jeopardy("localhost:9090");
	        
	        // Create and set up the content pane
	        JeopardyPanel gamePanel = new JeopardyPanel(game);
	        gamePanel.setOpaque(true);
	        frame.setContentPane(gamePanel);
	        frame.pack();
	        frame.setVisible(true);
	    }

}
