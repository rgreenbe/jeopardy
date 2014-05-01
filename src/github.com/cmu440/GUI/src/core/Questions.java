package core;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.Scanner;

public class Questions {
	private ArrayList<Question> q;
	private Scanner s;
	/*
	 * Questions gathers all questions from the .txt file and turns them into 
	 * question objects that have value options, answers
	 */
	public Questions() throws FileNotFoundException{
		q=new ArrayList<Question>();
		String file=("assets/Questions.txt");
		String question;
		int i =0;
		s = new Scanner (new File(file));
		while(s.hasNextLine()){
			question=s.nextLine();
			if (question.startsWith("Question"+Integer.toString(i)))
			{
				question=question.replace("Question"+Integer.toString(i)+":","");
				q.add(new Question(question));
			}else if(question.startsWith("Options"+Integer.toString(i))){
				question=question.replace("Options"+Integer.toString(i)+":","");
				q.get(i).setOptions(question);
			}else if(question.startsWith("Value"+Integer.toString(i))){
					question=question.replace("Value"+Integer.toString(i)+":","");
					q.get(i).setValue(Integer.parseInt(question));
			}else if(question.startsWith("Answer"+Integer.toString(i))){
				question=question.replace("Answer"+Integer.toString(i)+":","");
				q.get(i).setAnswer(Integer.parseInt(question));
				i+=1;
				
			}

		}
	}
	public  ArrayList<Question> questions(){
		return q;
	}
}

