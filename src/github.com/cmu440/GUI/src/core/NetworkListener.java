package core;

import java.io.IOException;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.UnknownHostException;

public class NetworkListener implements Runnable {
	private ServerSocket s;
	private InputStream in;
	private Jeopardy j;
	private final String join, buzz, answer, question;

	public NetworkListener(ServerSocket readSocket, Jeopardy j) {
		System.out.println("Starting network listener");
		this.s = readSocket;
		this.j = j;
		this.join = "{\"JoinRep\":";
		this.buzz = "{\"Buzz\":";
		this.answer = "{\"Answer\":";
		this.question = "{\"Question\":";
	}

	private String jsonString(String command, String json) {
		String withoutCommand = (json.replace(command, ""));
		withoutCommand = withoutCommand.replace("}\n", "");
		return withoutCommand.trim().replaceAll("\0", "");
	}

	@Override
	public void run() {
		while (true) {
			Socket socket = null;
			try {
				socket = s.accept();
				in = socket.getInputStream();
			} catch (IOException e) {
				e.printStackTrace();
			}
			while (socket.isConnected()) {
				byte[] buf = new byte[1000];
				try {
					in.read(buf);
					String msg = new String(buf, "UTF-8");
					if (msg.startsWith(join)) {
						j.joined(jsonString(join, msg));
					} else if (msg.startsWith(buzz)) {
						System.out.println("MESSAGE: ");
						System.out.println(msg);
						System.out.println("DONE");
						j.buzzed(jsonString(buzz, msg));
					} else if (msg.startsWith(answer)) {
						j.answered(jsonString(answer, msg));
					} else if (msg.startsWith(question)) {
						j.selectedQuestion(jsonString(question, msg));
					}

				} catch (IOException e) {
					e.printStackTrace();
				} catch (InterruptedException e) {
					e.printStackTrace();
				}

			}

		}
		
	}

}
