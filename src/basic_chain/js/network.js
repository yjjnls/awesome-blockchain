'use strict';

var net = require("net");
var Msg = require("./message");
var EventEmitter = require('events').EventEmitter;
var Promise = require("bluebird");

var PORT = 8000;

class Node extends EventEmitter {
    constructor(id) {
        super();

        this.id_ = id;
        this.peers_ = {};

        let self = this;
        this.server_ = net.createServer((socket) => {
            socket.setEncoding('utf8');
            socket.on('data', (data) => { self.on_data(data, socket); });
            socket.on('end', () => { self.remove_peer(socket); });
        });
        this.server_.listen(PORT + id);
    }
    async start() {
        for (var i = 0; i < 5; ++i) {
            var remote_id = Math.floor(Math.random() * 20); // [0, 20)
            if (remote_id !== this.id_ && !this.peers_[remote_id]) {
                let self = this;
                // console.log(`${this.id_}-->${remote_id}`);
                var socket = net.createConnection({ port: (PORT + remote_id) });
                await new Promise((resolve, reject) => {
                    socket.on('connect', () => {
                        resolve();
                    });
                    socket.on('error', function (e) {
                        resolve();
                    });
                    socket.setEncoding('utf8');
                    socket.on('data', (data) => { self.on_data(data, socket); });
                });
                // console.log(`id: ${self.id_} connected to remote_id: ${remote_id}`);
                let data = Msg.connection(self.id_);
                self.send(socket, data);
                self.add_peer(socket, remote_id);

            }
        }
    }
    on_data(data, socket) {
        try {
            var arr = data.split("\r\n");
            for (var i = 0; i < arr.length; ++i) {
                if (arr[i] == '') continue;
                let obj = JSON.parse(arr[i]);
                if (obj.type == Msg.type.Connection) {
                    // if data is connection info, add peer and response
                    let remote_id = obj.data;
                    this.add_peer(socket, remote_id);
                    // console.log(`node ${this.id_} receive connection: ${remote_id}`);
                    // console.log(`${this.id_}-->${remote_id}`);
                } else {
                    // else emit msg to blockchain
                    this.emit("message", obj);
                }
            }
        } catch (err) {
            console.log("=========================");
            console.log(`node: ${this.id_}\t receive msg error`);
            console.log(err);
            console.log(err.message);
            console.log(data);
            console.log(arr.length);
            console.log("=========================");
            throw new Error();
        }

    }
    send(socket, data) {
        if (typeof socket === 'number') {
            socket = this.peers_[socket];
        }
        if (typeof data === 'object') {
            data = JSON.stringify(data);
        }
        socket.write(data + "\r\n");
    }
    broadcast(data) {
        for (var index in this.peers_) {
            let socket = this.peers_[index];
            this.send(socket, data);
        }
    }
    add_peer(socket, remote_id) {
        if (!this.peers_[remote_id]) {
            this.peers_[remote_id] = socket;
            // console.log(`${this.id_}-->${remote_id}`);
        }
    }
    remove_peer(socket) {
        for (var index in this.peers_) {
            if (this.peers_[index] == socket) {
                delete this.peers_[index];
                break;
            }
        }
    }
    list_peers() {
        let peer_ids = [];
        for (var index in this.peers_) {
            peer_ids.push(index);
        }
        return peer_ids;
    }
}

module.exports = Node;