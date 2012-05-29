# coding: utf-8
"""Connection to hatcogd"""

from __future__ import unicode_literals


class Server(object):
    """Local proxy for remote hatcogd server"""

    def __init__(self, sock):

        self.conn = sock

    def write(self, msg):
        """Send a string message to the server"""
        if not msg:
            return
        msg += "\n"
        self.conn.sendall(msg.encode("utf8"))

    def stop(self):
        """Close server connection"""
        self.conn.close()

    def receive_one(self):
        """Listen for data on conn, return it on queue"""

        data = []

        char = self.conn.recv(1)
        while char != "\n":
            data.append(char)
            char = self.conn.recv(1)

        received = "".join(data)
        return received.decode("utf8")
