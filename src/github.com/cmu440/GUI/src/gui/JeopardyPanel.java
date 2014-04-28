package gui;

import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Dimension;
import java.awt.GridLayout;
import java.util.List;

import javax.swing.JButton;
import javax.swing.JLabel;
import javax.swing.JPanel;
import javax.swing.border.EmptyBorder;
import javax.swing.border.LineBorder;

import core.GameChangeListener;
import core.GameInfo;
import core.Jeopardy;
import core.Question;

public class JeopardyPanel extends JPanel implements GameChangeListener{
	private final Jeopardy j;
	private final int WIDTH,LENGTH,BWIDTH,BLENGTH;
	private JButton[][] questions;
	private  GameInfo gameInfo;
    private JPanel headPanel,footPanel,gridPanel,sidePanel;
	  public JeopardyPanel(Jeopardy game) {
		 WIDTH=880;
		 LENGTH=700;
		 BWIDTH=500;
		 BLENGTH=600;
		 j=game;
		 gameInfo=game.Info();
		 headPanel=createHeadPanel();
		 gridPanel=createGrid();
		 footPanel=createFootPanel();
		 sidePanel=createSidePanel();
		 game.addListener(this);
		 initGui();
	  }

	private JPanel createSidePanel() {
        JPanel side=new JPanel();
        side.setOpaque(true);
        side.setPreferredSize(new Dimension(200,LENGTH));
        side.setLayout(new GridLayout(10,1));
        JButton join=new JButton();
        join.addActionListener(new JoinListener(j));
        join.setText("Join Game");
        JPanel pad=new JPanel();
        pad.setBorder(new EmptyBorder(5,5,5,5));
        pad.add(join);      
        side.add(pad);
        JLabel players=new JLabel("Players",JLabel.CENTER);
        side.add(players);
		return side;
	}

	private JPanel createGrid() {
		setBorder(new EmptyBorder(30,30,30,30));
		JPanel grid=new JPanel();
		grid.setOpaque(true);
		setPreferredSize(new Dimension(WIDTH,LENGTH));
		List<List<Integer>> board=gameInfo.board();
		int padding=5;
		int rows=board.size();
		int cols=board.get(0).size();
		grid.setMaximumSize(new Dimension(BWIDTH,BLENGTH));
		grid.setBorder(new EmptyBorder(30,30,30,30));
		grid.setLayout(new GridLayout(rows+1,cols));
		addCategories(grid);
		questions=new JButton[rows][cols];
		for (int row=0;row<rows;row++){
			for(int col=0;col<cols;col++){
				int val= board.get(row).get(col);
				JPanel pad=new JPanel();
				pad.setOpaque(false);
				pad.setBorder(new EmptyBorder(10,10,10,10));
				JButton q=new JButton();
		
				q.setOpaque(true);
				q.addActionListener(new ChooseQuestion(row,col,q,j));
				pad.add(q);
				q.setPreferredSize(new Dimension(50,50));
				grid.add(pad);
				q.setText(Integer.toString(val));
			

			}
		}
		return grid;
	}

	private void addCategories(JPanel grid) {
		List<String> categories =gameInfo.Categories();
		for(int i=0;i<categories.size();i++){
			JLabel category=new JLabel(categories.get(i),JLabel.CENTER);
			grid.add(category);
		}	
	}
	private void initGui() {
		setLayout(new BorderLayout());
        Color darkCyan = new Color(10, 135, 171);
        setPreferredSize(new Dimension(WIDTH,LENGTH));
        
        add(headPanel,BorderLayout.NORTH);
        add(gridPanel,BorderLayout.CENTER);
        add(sidePanel,BorderLayout.EAST);
        add(footPanel,BorderLayout.SOUTH);
        setBackground(darkCyan);
        setOpaque(true);

			
	}
    private JPanel createFootPanel() {
        JPanel panel = new JPanel();
        JLabel copyright = new JLabel("Copyright © 2014 Gode Greenberg 15-440 Jeopardy!",JLabel.CENTER);
        panel.add(copyright);
        return panel;
    }
    private JPanel createHeadPanel() {
        JPanel panel = new JPanel();
        panel.setLayout(new BorderLayout());
        
        
        JLabel copyright = new JLabel("Welcome To Jeopardy! Review for the 15440 Final");
        panel.add(copyright,BorderLayout.NORTH);
        return panel;
    }

	@Override
	public void selectQuestion(Question q,int row,int col) {
		System.out.println("Selecting question: "+q.question());
	    gridPanel.removeAll();

        gridPanel.setLayout(new GridLayout(5,1));
        JLabel question=new JLabel(q.question());
        gridPanel.add(question);
        gridPanel.add(new JLabel("Select one of the options below"));
        for(int i=0;i<q.options().size();i++){
        	JButton option=new JButton();
        	option.setText(q.options().get(i));
        	gridPanel.add(option);
        }
        gridPanel.validate();
        gridPanel.repaint();
        
		
	}

}
