var express = require('express')
var WebSocket = require('ws')
var ShareDB = require('sharedb')
var WebSocketJSONStream = require('@teamwork/websocket-json-stream')

var app = express()
var server = http.createServer(app)
var webSocketServer = WebSocket.Server({server: server})

var backend = new ShareDB()
webSocketServer.on('connection', (webSocket) => {
  var stream = new WebSocketJSONStream(webSocket)
  backend.listen(stream)
})

server.listen(8080)