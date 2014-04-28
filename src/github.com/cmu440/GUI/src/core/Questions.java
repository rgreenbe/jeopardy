package core;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.Scanner;

public class Questions {
	private ArrayList<Question> q;
	public Questions() throws FileNotFoundException{
		q=new ArrayList<Question>();
		String file=("assets/Questions.txt");
		Scanner s;
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
			}else if(question.startsWith("Answer"+Integer.toString(i))){
				question=question.replace("Answer"+Integer.toString(i)+":","");
				q.get(i).setAnswer(Integer.parseInt(question));
				i+=1;
				
			}

		}
		System.out.println(q.size());
	
	}
	public  ArrayList<Question> questions(){
		return q;
	}
}

