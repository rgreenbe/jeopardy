package gui;

import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Dimension;
import java.awt.Font;
import java.awt.GridLayout;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import javax.imageio.ImageIO;
import javax.swing.BorderFactory;
import javax.swing.ImageIcon;
import javax.swing.JButton;
import javax.swing.JLabel;
import javax.swing.JPanel;
import javax.swing.border.EmptyBorder;
import javax.swing.border.LineBorder;

import core.GameChangeListener;
import core.GameInfo;
import core.Jeopardy;
import core.Question;

public class JeopardyPanel extends JPanel implements GameChangeListener {
	private final Jeopardy j;
	private final int WIDTH, LENGTH, BWIDTH, BLENGTH, MAXGUESSES;
	private JButton[][] questions;
	private JButton[] options;
	private ArrayList<JLabel> playerScores;
	private GameInfo gameInfo;
	private JButton buzzer, join;
	private JLabel result;
	private ArrayList<Integer> playersWhoGuessed;
	private JPanel headPanel, footPanel, gridPanel, sidePanel;
	private int rows, cols;

	public JeopardyPanel(Jeopardy game) {
		MAXGUESSES = 3;
		playersWhoGuessed = new ArrayList<Integer>();
		playerScores = new ArrayList<JLabel>();
		result = new JLabel();
		WIDTH = 880;
		LENGTH = 700;
		BWIDTH = 500;
		BLENGTH = 600;
		rows = game.Info().board().size();
		cols = game.Info().board().get(0).size();
		j = game;
		buzzer = new JButton();
		buzzer.setText("Buzz");
		buzzer.setBackground(Color.RED);
		buzzer.setOpaque(true);
		buzzer.addActionListener(new Buzzer(j));
		gameInfo = game.Info();
		headPanel = createHeadPanel();
		initBoardQuestions(rows, cols);
		gridPanel = new JPanel();
		createGrid();
		footPanel = createFootPanel();
		footPanel.setOpaque(false);
		sidePanel = createSidePanel();
		game.addListener(this);
		initGui();
	}

	private JPanel createSidePanel() {
		Font f = new Font("Helvetica", Font.PLAIN, 25);
		JPanel side = new JPanel();
		side.setOpaque(true);
		side.setPreferredSize(new Dimension(200, LENGTH));
		side.setLayout(new GridLayout(10, 1));
		join = new JButton();
		join.addActionListener(new JoinListener(j, join));
		join.setFont(f);
		join.setText("Join Game");
		JPanel pad = new JPanel();
		pad.setBorder(new EmptyBorder(5, 5, 5, 5));
		pad.add(join);
		side.add(pad);
		JLabel players = new JLabel("Players", JLabel.CENTER);
		side.add(players);
		return side;
	}

	private void createGrid() {
		setBorder(new EmptyBorder(30, 30, 30, 30));
		;
		gridPanel.setOpaque(true);
		setPreferredSize(new Dimension(WIDTH, LENGTH));
		List<List<Integer>> board = gameInfo.board();
		int padding = 5;
		gridPanel.setMaximumSize(new Dimension(BWIDTH, BLENGTH));
		gridPanel.setBorder(BorderFactory.createLineBorder(Color.black, 3));
		gridPanel.setLayout(new GridLayout(rows + 1, cols));
		Color jeopardyBlue = new Color(6, 12, 233);
		addCategories(gridPanel);
		Font f = new Font("Helvetica", Font.PLAIN, 25);
		for (int row = 0; row < rows; row++) {
			for (int col = 0; col < cols; col++) {
				int val = board.get(row).get(col);
				JPanel pad = new JPanel();
				pad.setBorder(BorderFactory.createLineBorder(Color.black, 3));
				pad.setBackground(jeopardyBlue);
				JButton q = questions[row][col];
				q.setOpaque(false);
				q.setForeground(Color.blue);
				q.setEnabled(false);
				pad.add(q);
				q.setPreferredSize(new Dimension(100, 100));
				gridPanel.add(pad);
				//q.setText(Integer.toString(val));
				q.setFont(f);

			}
		}

	}

	private void initBoardQuestions(int rows, int cols) {
		List<List<Integer>> board = gameInfo.board();
		ArrayList<Question> qFromBoard = j.questions();
		int index;
		questions = new JButton[rows][cols];
		for (int row = 0; row < rows; row++) {
			for (int col = 0; col < cols; col++) {
				JButton q = new JButton();
				questions[row][col] = q;
				index=row*cols+col;
				q.setText(Integer.toString(qFromBoard.get(index).value()));
				q.addActionListener(new ChooseQuestion(row, col, q, j));

			}
		}

	}

	private void addCategories(JPanel grid) {
		Color jeopardyBlue = new Color(6, 12, 233);
		Font f = new Font("Helvetica", Font.BOLD, 25);
		List<String> categories = gameInfo.Categories();
		for (int i = 0; i < categories.size(); i++) {
			JLabel category = new JLabel(categories.get(i), JLabel.CENTER);
			category.setBackground(jeopardyBlue);
			category.setForeground(Color.white);
			category.setOpaque(true);
			category.setBorder(BorderFactory.createLineBorder(Color.black, 3));
			category.setFont(f);
			grid.add(category);
		}
	}

	private void initGui() {
		Color jeopardyBlue = new Color(6, 12, 233);
		setLayout(new BorderLayout());
		setPreferredSize(new Dimension(WIDTH, LENGTH));
		add(headPanel, BorderLayout.NORTH);
		add(gridPanel, BorderLayout.CENTER);
		add(sidePanel, BorderLayout.EAST);
		add(footPanel, BorderLayout.SOUTH);
		setBackground(jeopardyBlue);
		setOpaque(true);

	}

	private JPanel createFootPanel() {
		JPanel panel = new JPanel();
		JLabel copyright = new JLabel(
				"Copyright � 2014 Gode Greenberg 15-440 Jeopardy!",
				JLabel.CENTER);
		panel.add(copyright);
		copyright.setForeground(Color.white);
		return panel;
	}

	private JPanel createHeadPanel() {
		JPanel panel = new JPanel();
		panel.setLayout(new BorderLayout());
		Font f = new Font("Helvetica", Font.BOLD, 25);
		JLabel header = new JLabel("Jeopardy! 15440 Review", JLabel.CENTER);
		header.setFont(f);
		header.setForeground(Color.white);
		panel.add(header, BorderLayout.NORTH);
		panel.setOpaque(false);
		return panel;
	}

	@Override
	public void selectQuestion(Question q, int row, int col) {
		System.out.println("Selecting question: " + q.question());
		questions[row][col].setEnabled(false);
		gridPanel.removeAll();
		gridPanel.add(new JLabel("Click Buzz when you know the answer"));
		gridPanel.setLayout(new GridLayout(10, 1));
		JLabel question = new JLabel(q.question());
		JPanel pad = new JPanel();
		pad.setOpaque(false);
		pad.setBorder(new EmptyBorder(10, 10, 10, 10));
		gridPanel.add(question);
		pad.add(buzzer);
		gridPanel.add(pad);
		gridPanel.add(result);
		gridPanel.add(new JLabel("Select one of the options below"));
		options = new JButton[q.options().size()];
		for (int i = 0; i < q.options().size(); i++) {
			pad = new JPanel();
			pad.setOpaque(false);
			pad.setBorder(new EmptyBorder(10, 10, 10, 10));
			JButton option = new JButton();
			option.addActionListener(new ChooseAnswer(i, j));
			option.setText(q.options().get(i));
			options[i] = option;
			option.setEnabled(false);
			pad.add(option);
			gridPanel.add(pad);
		}
		gridPanel.validate();
		gridPanel.repaint();

	}

	@Override
	public void buzzedIn(int playerID) {
		buzzer.setText("Player: " + Integer.toString(playerID) + " Buzzed");
		buzzer.setEnabled(false);
		if (j.playerID() == playerID) {
			for (int i = 0; i < options.length; i++) {
				options[i].setEnabled(true);
			}
		}

	}

	private void resetGrid() throws InterruptedException {
		Thread.sleep(1400);
		gridPanel.removeAll();
		createGrid();
		gridPanel.validate();
		gridPanel.repaint();
		buzzer.setEnabled(true);
		buzzer.setText("Buzz");
		enableQuestions();

	}

	@Override
	public void answeredQuestion(int playerID, int score, int choice)
			throws InterruptedException {
		System.out.println("Score for player :"+playerID+" is "+score);
		options[choice].setOpaque(true);
		options[choice].setEnabled(false);

		for (int i = 0; i < options.length; i++) {
			options[i].setEnabled(false);

		}
		if (j.currentQuestion().answer() == choice) {
			System.out.println("RIGHT");
			result.setText("RIGHT ANSWER");
			options[choice].setBackground(Color.GREEN);
			resetGrid();
			playersWhoGuessed.removeAll(playersWhoGuessed);
		} else {
			result.setText("Wrong Answer");
			playersWhoGuessed.add(playerID);
			options[choice].setBackground(Color.RED);
			buzzer.setText("Buzz");
			if (MAXGUESSES > playersWhoGuessed.size()) {
				if (!playersWhoGuessed.contains(j.playerID())) {
					buzzer.setEnabled(true);
				}

			} else {
				playersWhoGuessed.removeAll(playersWhoGuessed);
				resetGrid();
			}
		}
		playerScores.get(playerID).setText(
				"Player score: " + Integer.toString(score));

	}

	@Override
	public void startGame() {
		sidePanel.setLayout(new GridLayout(5, 1));
		System.out.println("adding players ");
		for (int i = 0; i < j.players().size(); i++) {
			BufferedImage image;
			try {
				image = ImageIO.read(new File("assets/gopher_player"
						+ Integer.toString(i + 1) + ".png"));
				image.getScaledInstance(image.getWidth(), image.getHeight(), 1);
				JLabel playerIcon = new JLabel(new ImageIcon(image));

				playerIcon.setPreferredSize(new Dimension(80, 80));

				JPanel pbox = new JPanel();
				pbox.setOpaque(true);
				pbox.setLayout(new GridLayout(4, 1));
				pbox.add(new JLabel("Player: " + Integer.toString(i)));
				pbox.add(playerIcon);
				JLabel score = new JLabel("Score: 0");
				pbox.add(score);
				playerScores.add(score);

				if (i == j.playerID()) {
					pbox.add(new JLabel("Your Player"));
				} else {
					pbox.add(new JLabel("Opponent"));
				}
				sidePanel.add(pbox);
				sidePanel.validate();
				sidePanel.repaint();
			} catch (IOException ex) {

			}
		}
		enableQuestions();

	}

	private void enableQuestions() {
		System.out.println("current player is " + j.currentPlayer()
				+ "and player ID " + j.playerID());
		Boolean[][] alreadyAnswered = j.SelectedQuestions();
		for (int row = 0; row < rows; row++) {
			for (int col = 0; col < cols; col++) {
				if (j.playerID() == j.currentPlayer()) {
					System.out.println("Please enable");
					if (alreadyAnswered[row][col]) {
						questions[row][col].setEnabled(false);
					} else {

						questions[row][col].setEnabled(true);
					}
				} else {
					questions[row][col].setEnabled(false);
				}
			}

		}
	}
}
