package core;

import gui.JeopardyPanel;

import javax.swing.JFrame;
import javax.swing.SwingUtilities;

public class Main {
	  private static final String NAME = "Scrabble";

	    public static void main(String[] args) {
	        SwingUtilities.invokeLater(new Runnable() {
	            @Override
	            public void run() {
	                createAndShowJeopardy();
	            }
	        });
	    }
	    private static void createAndShowJeopardy() {
	        // Create and set-up the window.
	        JFrame frame = new JFrame(NAME);
	        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);

	        Jeopardy game = new Jeopardy("{PlayerID:3,Test:[true,true],"
	        		+ "board:[[200,200,200,200,200],"
	        		+ "[400,400,400,400,400],"
	        		+ "[600,600,600,600,600],"
	        		+ "[800,800,800,800,800]],categories:[poop,hello,world,mgode,rgreenbe]"
	        		+ "}");
	        
	        // Create and set up the content pane
	        JeopardyPanel gamePanel = new JeopardyPanel(game);
	        gamePanel.setOpaque(true);
	        frame.setContentPane(gamePanel);
	        frame.pack();
	        frame.setVisible(true);
	    }

}