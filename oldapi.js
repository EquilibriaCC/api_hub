var express = require('express');
var app = express();
var fs = require("fs");
const axios = require('axios');
const {HLTV} = require('hltv')
const cheerio = require('cheerio')
var cors = require('cors');
var fetch = require('isomorphic-unfetch');

app.use(cors());

// Counts Pages
let page_visits = {}
let visits = async function (req, res, next) {
    async function count() {
        return new Promise(resolve => {
            let counter = page_visits[req.originalUrl.toLowerCase()];
            if (counter || counter === 0)
                page_visits[req.originalUrl.toLowerCase()] = counter + 1;
            else
                page_visits[req.originalUrl.toLowerCase()] = 1;
            resolve([req.originalUrl, page_visits[req.originalUrl.toLowerCase()]])
        })
    }

    let counter = await count()
    next();
};
app.use(visits);

app.get('/api/v1/number_of_requests', function (req, res) {
    fs.readFile(__dirname + "/" + "number_of_requests.json", 'utf8', function (err, data) {
        res.end(data)
    })
})

app.get('/api/v1/emission', function (req, res) {
    fs.readFile(__dirname + "/" + "emissions.json", 'utf8', function (err, data) {
        res.end(data)
    })
})

app.get('/api/v1/number_of_nodes', function (req, res) {
    fs.readFile(__dirname + "/" + "number_of_nodes.json", 'utf8', function (err, data) {
        res.end(data)
    })
})

app.get('/api/v1/matches', function (req, res) {
    fs.readFile(__dirname + "/" + "matches.json", 'utf8', function (err, data) {
        res.end(data)
    })
})

app.get('/api/v1/metals', function (req, res) {
    fs.readFile(__dirname + "/" + "metals.json", 'utf8', function (err, data) {
        res.end(data)
    })
})
app.get('/api/v1/explorer', function (req, res) {
    fs.readFile(__dirname + "/" + "explorer.json", 'utf8', function (err, data) {
        res.end(data)
    })
})

app.get('/api/v1/service_nodes', function (req, res) {
    fs.readFile(__dirname + "/" + "service_nodes.json", 'utf8', function (err, data) {
        res.end(data)
    })
})
app.get('/api/v1/explorer/transactions/:id/:start_height', function (req, res) {
    fs.readFile(__dirname + "/" + "explorer.json", 'utf8', function (err, data) {
        let height = JSON.parse(data)['height'] - 1 - req.params.start_height

        async function transactions(number_of_transaction, height) {
            let txs = []

            async function pull(height) {
                return new Promise(resolve => {
                    async function call(height) {
                        return new Promise(call_promise => {
                            fetch("http://sanfran.equilibria.network:9231/json_rpc", {
                                body: JSON.stringify({
                                    "jsonrpc": "2.0", "id": "0", "method": "get_block", "params": {"height": height}
                                }),
                                headers: {"Content-Type": "application/json"},
                                method: "POST"
                            }).then(res => res.json()).then((result) => {
                                    txs.push(result.result.miner_tx_hash)
                                    for (let item in (result.result.tx_hashes)) txs.push(result.result.tx_hashes[item])
                                    call_promise()
                                },
                                (error) => {
                                    console.log(error)
                                }
                            )
                        })
                    }

                    let promises = []
                    while (number_of_transaction >= 0) {
                        promises.push(call(height))
                        number_of_transaction--
                        height--
                    }
                    Promise.all(promises).then(() => {
                        resolve(txs)
                    })
                })
            }

            let tx_list = await pull(height)
            let transaction_list = {}
            fetch("http://sanfran.equilibria.network:9231/get_transactions", {
                body: JSON.stringify({"txs_hashes": tx_list, "decode_as_json": true}),
                headers: {"Content-Type": "application/json"},
                method: "POST"
            }).then(res => res.json())
                .then(
                    (result) => {
                        let height_txs = []
                        for (let item in result.txs) {
                            let transaction = JSON.parse(result.txs[item].as_json)
                            let total_output = 0
                            for (let output in transaction['vout']) {
                                total_output = total_output + transaction['vout'][output]['amount'] / 10000
                            }
                            let obj = {
                                "block": result.txs[item].block_height, "tx_hash": result.txs[item].tx_hash,
                                "unlock_block": transaction['unlock_time'], "outputs": transaction['vout'],
                                "inputs": transaction['vin'], "timestamp": result.txs[item].block_timestamp,
                                "total_output": total_output,
                            }
                            if (transaction['rct_signatures']) {
                                obj = {...obj, ...{"ringCT_type": transaction['rct_signatures']['type']}}
                                if (transaction['rct_signatures']['txnFee'])
                                    obj = {...obj, ...{"txnFee": transaction['rct_signatures']['txnFee']}}
                            }
                            height_txs.push(obj)
                        }
                        for (let hash in height_txs) {
                            if (!transaction_list[height_txs[hash].block]) {
                                transaction_list[height_txs[hash].block] = [height_txs[hash]]
                            } else {
                                transaction_list[height_txs[hash].block].push(height_txs[hash])
                            }
                        }
                        let new_list = []
                        for (let block in transaction_list) new_list.push({[block]: transaction_list[block]})
                        res.end(JSON.stringify(new_list.reverse(), null, 4))
                    },
                    (error) => {
                        console.log(error)
                    }
                )
        }

        transactions(req.params.id, height)
    })
})
app.get('/api/v1/explorer/txpool', function (req, res) {
    fetch("http://sanfran.equilibria.network:9231/get_transaction_pool", {
        headers: {"Content-Type": "application/json"},
        method: "POST"
    }).then(res => res.json()).then((result) => {
            let transactions = []
            for (let item in result.transactions) {
                let transaction_json = JSON.parse(result.transactions[item].tx_json)
                let transaction = result.transactions[item]
                let obj = {
                    "hash": transaction['id_hash'],
                    "timestamp": transaction['receive_time'],
                    "version": transaction_json['version'],
                    "outputs": transaction_json.vout,
                    "inputs": transaction_json.vin,
                }
                if (transaction_json['rct_signatures']) {
                    obj = {...obj, ...{"ringCT_type": transaction_json['rct_signatures']['type']}}
                    if (transaction_json['rct_signatures']['txnFee']) {
                        obj = {...obj, ...{"fee": transaction_json['rct_signatures']['txnFee']}}
                    }
                }
                transactions.push(obj)
            }
            res.end(JSON.stringify(transactions, null, 4))
        },
        (error) => {
            console.log(error)
        }
    )
})
app.get('/api/v1/explorer/tx/:id', function (req, res) {
    fs.readFile(__dirname + "/" + "explorer.json", 'utf8', function (err, data) {
        let current_height = JSON.parse(data)['height']
        fetch("http://sanfran.equilibria.network:9231/get_transactions", {
            body: JSON.stringify({"txs_hashes": [req.params.id], "decode_as_json": true}),
            headers: {"Content-Type": "application/json"},
            method: "POST"
        }).then(res => res.json()).then((result) => {
                for (let item in result.txs) {
                    let transaction = JSON.parse(result.txs[item].as_json)
                    let total_output = 0
                    for (let output in transaction['vout']) {
                        total_output = total_output + transaction['vout'][output]['amount']
                    }
                    let obj = {
                        "block": result.txs[item].block_height,
                        "tx_hash": result.txs[item].tx_hash,
                        "unlock_block": transaction['unlock_time'],
                        "version": transaction['version'],
                        "outputs": transaction['vout'],
                        "inputs": transaction['vin'],
                        "timestamp": result.txs[item].block_timestamp,
                        "total_output": total_output,
                        "ringCT_type": transaction['rct_signatures']['type'],
                        "confirmation": current_height - result.txs[item].block_height
                    }
                    if (transaction['rct_signatures']['txnFee']) {
                        obj = {...obj, ...{"txnFee": transaction['rct_signatures']['txnFee']}}
                    }
                    res.end(JSON.stringify(obj, null, 4))
                }

            },
            (error) => {
                console.log(error)
            }
        )
    })
})

app.get('/api/v1/match/results/:id', function (req, res) {
    fs.readFile(__dirname + "/" + "matches.json", 'utf8', function (err, data) {
        let matches = JSON.parse(data)
        for (let i = 0; i < matches.length; i++) {
            if (req.params.id === String(matches[i]['matchid'])) {
                let match = matches[i]
                res.end(JSON.stringify(match))
                break
            }
        }
    })
})
app.get('', function (req, res) {
    res.end(JSON.stringify('Go away bud'))
})


// Update the API Data
function updateMatches() {
    function update(jsonData) {
        let newJSON = []
        for (i = 0; i < jsonData.length; i++) {
            let tempData = {
                "matchid": jsonData[i]['id'],
                "team_1": jsonData[i]['team1'],
                "team_2": jsonData[i]['team2'],
                "Event": jsonData[i]['event'],
                "Completed": "false"
            }
            newJSON.push(tempData)
        }
        fs.writeFile("matches.json", JSON.stringify(newJSON, null, 4), function (err) {
            console.log('Updated Matches API on ' + Date())
            if (err) {
                console.log(err)
            }
        })
    }

    try {
        HLTV.getMatches().then((data) => {
            update(data);
        })
    } catch (err) {
        console.log("Failed updating matches." + err)
    }
}

function updatesMetals() {
    function update(data) {
        var prices = data['rates']

        fs.writeFile("metals.json", JSON.stringify(prices, null, 4), function (err) {
            console.log('Updated Metals API on ' + Date())
            if (err) {
                console.log(err)
            }
        })
    }

    try {
        axios.get('https://metals-api.com/api/latest?access_key=1k11a48r6tujxscqwyr3qkru3a370132o3i9co6d7tw71i5qws6p23ywy0jj&base=USD').then(response => {
            update(response.data)
        }).catch(err => {
            console.log('Metals API is down! ' + err)
        })
    } catch (err) {
        console.log("Failed updating Metals." + err)
    }
}

function updateExplorer() {
    function info(supply, marketCap, rank) {

        supply = supply['data']['coinbase'].toString()
        supply = supply.substring(0, 2) + "." + supply.substring(2, 3) + "M"

        let date = new Date()
        let day = date.getDate();
        let month = date.getMonth();
        let year = date.getFullYear();
        let hour = date.getHours()
        let minute = date.getMinutes()
        date = hour + ':' + minute + " " + day + "-" + (month + 1) + "-" + year;

        fetch("http://sanfran.equilibria.network:9231/json_rpc", {
            body: "{\"jsonrpc\":\"2.0\",\"id\":\"0\",\"method\":\"get_coinbase_tx_sum\",\"params\":{\"height\":370961,\"count\":1}}",
            headers: {
                "Content-Type": "application/json"
            },
            method: "POST"
        }).then(res => res.json())
            .then(
                (result) => {
                    fetch("http://sanfran.equilibria.network:9231/json_rpc", {
                        body: "{\"jsonrpc\":\"2.0\",\"id\":\"0\",\"method\":\"get_info\"}",
                        headers: {
                            "Content-Type": "application/json"
                        },
                        method: "POST"
                    })
                        .then(res => res.json())
                        .then(
                            (result) => {
                                // block height, media block size, difficulty, block time target, fork version
                                let block_size_medium = result.result.block_size_median
                                let difficulty = result.result.difficulty
                                let grey_peerlist_size = result.result.grey_peerlist_size
                                let height = result.result.height
                                let target_block_time = result.result.target
                                let tx_pool_size = result.result.tx_pool_size
                                let version = result.result.version
                                let white_peerlist_size = result.result.white_peerlist_size

                                let hash_rate = difficulty / 120
                                if (hash_rate > 1000000000) {
                                    hash_rate = (hash_rate / 1000000000).toLocaleString() + " GH/s"
                                } else if (hash_rate > 1000000) {
                                    hash_rate = (hash_rate / 1000000).toLocaleString() + " MH/s"
                                } else if (hash_rate > 1000) {
                                    hash_rate = (hash_rate / 1000).toLocaleString() + " kH/s"
                                }

                                fetch("http://sanfran.equilibria.network:9231/json_rpc", {
                                    body: "{\"jsonrpc\":\"2.0\",\"id\":\"0\",\"method\":\"" + "hard_fork_info" + "\"}",
                                    headers: {"Content-Type": "application/json"},
                                    method: "POST"
                                })
                                    .then(res => res.json())
                                    .then(
                                        (result) => {
                                            let hard_fork_height = result.earliest_height
                                            let fork_version = result.version
                                            fs.readFile(__dirname + "/" + "service_nodes.json", 'utf8', function (err, data) {
                                                let number_of_service_nodes = JSON.parse(data)[0].length
                                                let total_locked = 0
                                                for (let item in JSON.parse(data)[0]) {
                                                    total_locked = total_locked + Number(JSON.parse(data)[0][item]['staking_requirement'])
                                                }
                                                fetch("http://sanfran.equilibria.network:9231/json_rpc", {
                                                    body: JSON.stringify({
                                                        "jsonrpc": "2.0",
                                                        "id": "0",
                                                        "method": "get_staking_requirement",
                                                        "params": {"height": height - 1}
                                                    }),
                                                    headers: {"Content-Type": "application/json"},
                                                    method: "POST"
                                                }).then(res => res.json()).then((result) => {
                                                        let stakingreq = result.result.staking_requirement / 10000
                                                        fetch("http://sanfran.equilibria.network:9231/json_rpc", {
                                                            body: JSON.stringify({
                                                                "jsonrpc": "2.0",
                                                                "id": "0",
                                                                "method": "get_last_block_header",
                                                                "params": {"fill_pow_hash": true}
                                                            }),
                                                            headers: {"Content-Type": "application/json"},
                                                            method: "POST"
                                                        }).then(res => res.json()).then((result) => {
                                                                let blockreward = result.result.block_header.reward / 10000
                                                                let nodereward = Math.round((blockreward / 2) * 720 / number_of_service_nodes)

                                                                let data = {
                                                                    "numberofservicenodes": number_of_service_nodes,
                                                                    "hashrate": hash_rate,
                                                                    "difficulty": difficulty,
                                                                    "block_size_medium": block_size_medium,
                                                                    "height": height,
                                                                    "target_block_time": target_block_time,
                                                                    "grey_peerlist_size": grey_peerlist_size,
                                                                    "white_peerlist_size": white_peerlist_size,
                                                                    "version": version,
                                                                    "hard_fork_height": hard_fork_height,
                                                                    "hard_fork_version": fork_version,
                                                                    "tx_pool_size": tx_pool_size,
                                                                    "supply": supply,
                                                                    "stakingreq": Math.round(stakingreq),
                                                                    "blockreward": blockreward,
                                                                    "daily_emission": (blockreward * 720).toLocaleString(),
                                                                    "monthly_emission": (blockreward * 720 * 30).toLocaleString(),
                                                                    "annual_emission": (blockreward * 720 * 365).toLocaleString(),
                                                                    "annual_ROI": Math.round((nodereward / stakingreq) * 100 * 365) + '%',
                                                                    "total_locked": Math.round(total_locked / 10000).toLocaleString(),
                                                                    "nodereward": nodereward,
                                                                    "breakeven": Math.round(stakingreq / nodereward),
                                                                    "marketcap": marketCap,
                                                                    "rank": rank,
                                                                    "lastUpdate": date
                                                                }

                                                                fs.writeFile("explorer.json", JSON.stringify(data, null, 4), function (err) {
                                                                    console.log('Updated Local Explorer API on ' + Date())
                                                                    if (err) {
                                                                        console.log(err)
                                                                    }
                                                                })

                                                            },
                                                            (error) => {
                                                                console.log(error)
                                                            }
                                                        )
                                                    },
                                                    (error) => {
                                                        console.log(error)
                                                    }
                                                )
                                            })
                                        },
                                        (error) => {
                                            console.log(error)
                                        }
                                    )

                            },
                            (error) => {
                                console.log(error)
                            }
                        )

                },
                (error) => {
                    console.log(error)
                }
            )


    }

    function marketInfo(data) {

        let rank = data.market_cap_rank
        let marketCap = data.market_data.market_cap.usd

        axios.get('https://explorer.equilibria.network/api/emission').then(response => {
            info(response.data, marketCap, rank)
        })

    }

    axios.get("https://api.coingecko.com/api/v3/coins/triton?tickers=true&market_data=true&community_data=false&developer_data=false&sparkline=false").then(response => {
        marketInfo(response.data)
    }).catch(err => {
        console.log('Coin Gecko is down!')
    })


}

function serviceNodes() {
    let method = "get_service_nodes"
    fetch("http://sanfran.equilibria.network:9231/json_rpc", {
        body: "{\"jsonrpc\":\"2.0\",\"id\":\"0\",\"method\":\"" + method + "\"}",
        headers: {"Content-Type": "application/json"},
        method: "POST"
    })
        .then(res => res.json()).then((result) => {
            let service_nodes = result.result.service_node_states
            fetch("http://api.ili.bet/api/v1/explorer")
                .then(res => res.json()).then((result) => {
                    function quorum(service_nodes, height) {
                        fetch("http://sanfran.equilibria.network:9231/json_rpc", {
                            body: JSON.stringify({
                                "jsonrpc": "2.0", "id": "0", "method": "get_quorum_state", "params":
                                    {"height": height - 1}
                            }),
                            headers: {"Content-Type": "application/json"},
                            method: "POST"
                        }).then(res => res.json()).then((result) => {
                                let all_nodes = [service_nodes, result.result.nodes_to_test, result.result.quorum_nodes]
                                fs.writeFile("service_nodes.json", JSON.stringify(all_nodes, null, 4), function (err) {
                                    console.log('Updated Local Service Nodes API on ' + Date())
                                })
                            },
                            (error) => {
                                console.log(error)
                            }
                        )
                    }

                    quorum(service_nodes, result.height)
                },
                (error) => {
                    console.log(error)
                }
            )

        },
        (error) => {
            console.log(error)
        }
    )
}

function numberOfNodes() {

    fs.readFile(__dirname + "/" + "explorer.json", 'utf8', function (err, data) {
        fs.readFile(__dirname + "/" + "number_of_nodes.json", 'utf8', function (err, nodes) {
            if (err !== null) {
                nodes = [[JSON.parse(data).numberofservicenodes, Date.now()]]
            } else {
                nodes = JSON.parse(nodes)
                nodes.push([JSON.parse(data).numberofservicenodes, Date.now()])
                if (nodes.length > 60)
                    nodes = nodes.slice(nodes.length - 60)
            }
            fs.writeFile("number_of_nodes.json", JSON.stringify(nodes, null, 4), function (err) {
                console.log('Updated node number on ' + Date())
                if (err) {
                    console.log(err)
                }
            })
        })
    })

}

function getEmission() {
    fs.readFile(__dirname + "/" + "explorer.json", 'utf8', function (err, data) {
        let current_height = JSON.parse(data)['height']
        fs.readFile(__dirname + "/" + "emission.json", 'utf8', function (err, emissions) {
            if (err !== null) {
                fs.writeFile("emission.json", JSON.stringify([], null, 4), function (err) {
                    getEmission()
                    if (err) {
                        console.log(err)
                    }
                })
            } else {
                fetch("http://sanfran.equilibria.network:9231/json_rpc", {
                    body: JSON.stringify({
                        "jsonrpc": "2.0", "id": "0", "method": "get_coinbase_tx_sum", "params": {
                            "height": 0, "count": current_height - 1
                        }
                    }),
                    headers: {"Content-Type": "application/json"},
                    method: "POST"
                }).then(res => res.json()).then((result) => {
                        let supply = result.result.emission_amount / 10000
                        emissions = JSON.parse(emissions)
                        emissions.push([supply, Date.now()])
                        fs.writeFile("emissions.json", JSON.stringify(emissions, null, 4), function (err) {
                            console.log('Updated node number on ' + Date())
                            if (err) {
                                console.log(err)
                            }
                        })

                    },
                    (error) => {
                        console.log(error)
                    }
                )
            }
        })
    })
}

const matchesTime = 180 * 1000
const metalsTime = 60 * 1000 * 60 * 24
const numNodesTime = 60 * 1000 * 60 * 12
const explorerTime = 10 * 1000
const matchesInterval = setInterval(updateMatches, matchesTime)
const metalsInterval = setInterval(updatesMetals, metalsTime)
const explorerInterval = setInterval(updateExplorer, explorerTime)
const serviceNodesInterval = setInterval(serviceNodes, explorerTime)
const numberOfNodesInterval = setInterval(numberOfNodes, numNodesTime)
const emissionInterval = setInterval(getEmission, numNodesTime)

const server = app.listen(8080, function () {
    const host = server.address().address
    const port = server.address().port
    console.log("API now running on http://%s:%s", host, port)
})
