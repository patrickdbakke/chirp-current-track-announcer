const path = require("path");
const shelljs = require("shelljs");
const http = require("http");
const net = require("net");
const express = require("express");
const dgram = require("dgram");
const expect = require("chai").expect;
const currentPlaylist = require("./current_playlist");
const announcer = path.resolve(__dirname, "../announcer.go");

const IP_ADDRESS = "127.0.0.1";
const CHIRP_PORT = 8082;
const PROSTREAM_PORT = 8083;
const RDS_PORT = 8084;
const CHIRP_URL = `http://${IP_ADDRESS}:${CHIRP_PORT}/`;

const COMMAND = `go run ${announcer}`;
const PROSTREAM_COMMAND = `${COMMAND} --chirp ${CHIRP_URL} --prostream ${IP_ADDRESS} --port ${PROSTREAM_PORT} --runOnce true`;
const RDS_COMMAND = `${COMMAND} --chirp ${CHIRP_URL} --rds ${IP_ADDRESS} --rdsPort ${RDS_PORT} --runOnce true`;

const localTCPServer = (port, response) => {
  let server;
  return new Promise((resolve, reject) => {
    server = net.createServer(socket => {
      socket.on("data", data => {
        resolve(data.toString("utf8"));
      });
      if (response) {
        socket.write(response);
      }
      socket.end();
    });
    server.on("error", e => reject(e));
    server.listen(port);
  })
    .catch(e => {
      server.close();
      throw e;
    })
    .then(message => {
      server.close();
      return message;
    });
};

const localUDPServer = (port, response) => {
  let server;
  return new Promise((resolve, reject) => {
    server = dgram.createSocket("udp4");
    server.on("message", data => {
      resolve(data.toString("utf8"));
    });
    server.on("error", e => reject(e));
    server.bind(port);
  })
    .catch(e => {
      server.close();
      throw e;
    })
    .then(message => {
      server.close();
      return message;
    });
};

const run = command => {
  let logs = "";
  return new Promise(resolve => {
    const child = shelljs.exec(command, {
      silent: true,
      async: true
    });
    child.stdout.on("data", function(data) {
      logs += data.toString("utf8");
    });
    child.stdout.on("end", function() {
      resolve(logs);
    });
    child.stdin.end();
  });
};

describe("announcer", () => {
  let chirpServer, logs, message;

  beforeEach(() => {
    const chirp = express();
    chirp.get("/", function(req, res, next) {
      res.json(currentPlaylist);
      next();
    });
    chirpServer = chirp.listen(CHIRP_PORT);
  });

  afterEach(() => {
    chirpServer.close();
  });

  context("when it can connect to chirp", () => {
    context("when talking to rds", () => {
      const expectedMessage = `DPS='Go Away' by Jeff Parker on CHIRP Radio\n`;
      const expectedErrorLog = `The RDS Encoder did not like the input ${expectedMessage}`;
      context("when RDS can handle the data correctly", () => {
        beforeEach(() => {
          return Promise.all([
            localTCPServer(RDS_PORT, "Hello").then(m => (message = m)),
            run(RDS_COMMAND).then(l => (logs = l))
          ]);
        });

        it("should send the message to rds", () => {
          expect(message).to.equal(expectedMessage);
          expect(logs).to.equal("");
        });
      });

      context("when RDS cannot handle the data", () => {
        beforeEach(() => {
          return Promise.all([
            localTCPServer(RDS_PORT, "NO").then(m => (message = m)),
            run(RDS_COMMAND).then(l => (logs = l))
          ]);
        });

        it("should log the error", () => {
          expect(logs).to.include(expectedErrorLog);
        });
      });
    });

    context("when talking to prostream", () => {
      const expectedMessage = `t=Go Away - Jeff Parker | u=http://www.chirpradio.org\r\n`;
      const expectedErrorLog = `The RDS Encoder did not like the input ${expectedMessage}`;
      context("when prostream can handle the data correctly", () => {
        beforeEach(() => {
          return Promise.all([
            localUDPServer(PROSTREAM_PORT, "Hello").then(m => (message = m)),
            run(PROSTREAM_COMMAND).then(l => (logs = l))
          ]);
        });

        it("should send the message to prostream", () => {
          expect(message).to.equal(expectedMessage);
          expect(logs).to.equal("");
        });
      });
    });
  });
});
