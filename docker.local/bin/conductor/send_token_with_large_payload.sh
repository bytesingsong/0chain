#!/bin/bash

set -e
docker ps -a
printf '{"client_id":"1746b06bb09f55ee01b33b5e2e055d6cc7a900cb57c0a3a5eaabb8a0e7745802","client_key":"7b630ba670dac2f22d43c2399b70eff378689a53ee03ea20957bb7e73df016200fea410ba5102558b0c39617e5afd2c1843b161a1dedec15e1ab40543a78a518","keys":[{"public_key":"7b630ba670dac2f22d43c2399b70eff378689a53ee03ea20957bb7e73df016200fea410ba5102558b0c39617e5afd2c1843b161a1dedec15e1ab40543a78a518","private_key":"c06b6f6945ba02d5a3be86b8779deca63bb636ce7e46804a479c50e53c864915"}],"mnemonics":"cactus panther essence ability copper fox wise actual need cousin boat uncover ride diamond group jacket anchor current float rely tragic omit child payment","version":"1.0","date_created":"2021-08-04 18:53:56.949069945 +0100 BST m=+0.018986002"}' > ~/.zcn/testing.json;

./zwalletcli/zwallet --wallet testing.json faucet \
  --methodName pour --input "{Pay day}" --tokens 1;

./zwalletcli/zwallet send --to_client_id e7ebb698213b6bda097c0a14ccbe574356e99e9b666e4baeae540da1d9b51e7e --tokens .2 --desc "olj3mocKyyMyam6RtaRG7dwr7jgsvxUZo3KPmt36mv1QL3knItlAPlAdZi0KdNQ8VqgOyKjjvjsKzKuqRHTscCVu7o7IKBBKNoYe fj Qv6ISi9TbCnoBB d1RIE9fS1KDpNrv1lHsZP0 q779bG04rc6EfJhhorq3LFyL8rOo856gEHdJhczIGEK9shYCFO lQMGn4eT nypYw5sF4FRvslPqXYThETiF2jZQqlh9YDCXMtwpEUB2JgqkzkzKGlIjCC8Es4uF1N1kyri0eKDDuIwDlQIXBgzc VgGshihaMWd 0mLevHL9wYKbBgz8tI1wVKg0QxjQCBYUdeqUtI64ZNyiEiS6VkNM zyZLhxroscQ4PlOwIJM8F7VeZvz 7Cl5UO HIPLqu5AWDbF usrVH1CTl8Di3q AY o65vNlq17SJ2yVZNpcSLURLzybjL9vF2l vrvqoP7cHfDkDi4ZcCWtHQx6m zg4ZVOkkvWym05IhkwPPlr FGqM1FhrW94PzUTJXcOhukyfqhC73AKlNSuzvnBvJe8Ut6IzlmVReeg9PnXTXzMcG2EUDdJRFK4dP3HWvljRuQhq9F9wcisF18ls K6eLBPZ9iiCZGnBAEQJlbfRp GSNdkYoW HnYpLW77xfTH5TsQtayVKSiMqd7G0qs8EFdjMBSWTRXG5eEdYnlJNHbD6SA8MIEDbt sCuxSejQ1vYvTwfuV8 GbosOoa3k1V3RXae FxVmHRnys9EgFLwBxricfi91mKOKGZdqpy8iJ5 KpEdwQMY1POrYLPhO6HZYeLES fiYqkcC2a0UjB1l0nUhLUCmXGTwRExwH21DV4MEeVL x7IofcRrfWIvRqBIezAqn9WczBR39T2A9M7DJjOBi7P8fYMK0 4Sj2ECu c51Hns14QtGtvPK59iMmWquGT247YilTjGOZEgCw27XpQkuZkGp1nEDAERhwNIWilP6T6sWejF54 ZWpgEfqgRugTWkJzqTJPqvC 0tCt9kIY3UEfeXKvYjjcsES8S60QpiUlGCtbNJfT0xEXHaqvGrzrGRyKgPXzqgugk2OYlP4UsJDK52fLOXb68h9fu ABXiKZUawv v6 GbWZIdsIyX0A7 LUdZicpsBIBzXqCg2MeHBuQFcrJfmqjaPo3OYa1EITZZbXH83tEWXKyFnlQyz28mac4N9e4eIvk0edwHEPP Xs3qWvVc1UDtacAiWgBgLzyC2KBK3FR7V2g9mRUJJOTrpGBfCtAlkWK2EigkMmXlb69QI4bH6eL6b44tlgUX9cKZxpPlTLx7Kimd DW0xOdtG KR3pWUmO7wFXaWeEg2nAbvQadz9EWG5UUJPmjE9afzU7nAiu2oXGy6qelLzxdSOX01EQL02gWGODT6ELSNwB82uJoDs Cxe8MikJQhEqXtTmeEWpTfclT7ZSxOo1hlZEXWeTz1gkMbWukuEZOP3mRQ3DN80TnT05l 8w5bXffDoM2wn6RSQim1xD3FILaldN 947TzM1PjZ3P4e8lNHhVq4zuwlcsfiaiEsKGS6ZeSDI185rqopKa09Smu1WMPMNIwjlYIay wickpNhtOYdbc4pJi93TKujq1lTC SeXAEAr9kSQuQ0cIuV1cx83peCPIvwgqphwewVEYaEv5liepJ2T QFcrFCv ketJipgjaO3Um48ms0WTamJ7VVdQremprPGpFaCu AdXSZY7h cmcEW 1OrbWB21B89sWi3BM3JtD5MueEHgzIcQhHiIk0dKqPO86mIigodgr dxOUMD1YKzb7tAhIbY0CCq45bvXO1zC YyJ1YO 5jhvjNFW0blvCA5CQ zsthEnejBOhQ8i393JOwMZkqYz5JpODNiYfhnvlg2cAZYMSJcbaHSixvb8LUNa9zOk29znvubAS vxpBhyeAZBU2vQtjEFvyILpLOQlGYCi2WcTDK rOxt0cKyhZd1WiUuSBsQnIcaoBKnqgbkcDymAx89 LlfPOAqYGQcLvOGO9BRc8 EcMQjGfjMUp5CK9KDXKQMFovoEGI56XWPEIZkibGtkWYWfi4X1bBhDjHcqP1Q3cp9akEiBLzU1VDeufiD4ZoGTX1I2LhJVPYHQcg KkRj7 dRVADuXfr1jkUZJ6m17bzSRNohdCvJRwGZCmM 6ixwQkWqCDCg5EfL7LGXPvW Hkz6QCAKW28C9hfST2kTuDebtMXWWUFk  ZHlAaAiQv0PXF7M5 Zg6Lt ZW1CpT70yoWzOZRideevyjZKdc72Y7Yqawhdgvg9vH0bkhi5munryI7YagVPAK23Uzs7hVcO7XCS GwcFlSfopItfNg7QFKCAOcAyVajkVNVBEyXYqhZT2h0dCpxVqJ9cDgzpYH8Hff3  4uBZnJ5fqqDUfqdKDo0QYRsP12WBU3sCh0w p5jzacJb3DFvAO3sCXzAZw4SWDY21h8C0nLg8x4anpP4dR0Nf1JwSCZJgQp8DB8kKn3fK w09QTy3s5UL1JAE0SeHDb0S4GUnfh3 cMy6ta1JLhmmUmwucXp9Ic3cbwTpkL9YBt1Ha6hXw63kOq cYcVIVADHMLiwtvMoWQgFd5eMp674wd0d37BF 19kVo0AWwxXFAV0 xJYZwC2ZcMHkyke93Ls I LaL4SaY6Wk aJTyi6JLWjr6EH4UfXvZ Lr0GuzXAW447KohBYSGRTUxipPy91IaUjQmxyUdwDgk2a9 iRzpRHOkbMm4yGDle42azca4fnMqiB6NNNCWj5Gih24XF1v1HUmoTBVu8XYrARzwoXYzzQxJpTQAjk60jPqV7UBc2kI6jSQyveMd L7E3PW2aWHXeRiGsUMHore3ABsi4lybxrFCFJIUnN1EMUhpWP0LZQcm5YLQ9wtl4RA5J eGGnBu9BZtj0WBpSnaqQwSRtc8sV8c0 JtpUFHOEKIeHzfr6PgIlhMFmlgpjkxSoAY957hKadaoGvAsidBN5cVGTCYzgDv3YBDili4EkIimIFldIymAp730wjlDGKdVLxzZQ ymWTBRiPidtdPI1rdbCDDVkpEifsxBHVmpCVp93G7jbRhYHvLccof0br c2ilLsQyHPTf9ycym8EAYWmzl 4qL32YFWuvECVn5Jv 1KNDow9ju1hXPm HtWNT4DmIVePf0rBbpU9C6VVIBmT2dorChUjkwzrLQGiApOvkp6CUkydF0l5RX OHkJ4wT6d8EafNI2hMEr5s sGWV7u9dxNS3YZmBsH7HcTiU lC7eVbo3u2CtlbAzGUrf7dJXYesoXnQOflxf RD2Rzn47VHw3YZ77KK24XESvPjFNawvLX LlDO k7QZTSABja46KWq0YfkNxm9NaI4qK5eczr3r3qbkbW3jqmXxlvhHRPDJZg1yHeh4Ac9WNBLjJSjZPzwkqC7JhWUqBotHRifAXDwF VW PHw1Mi xdNJpjUtJXQFIHQVgeh6XjoDkVGEImwYIPfz3xw4Nqjbw6ofXBPZCvL513t0OjQRvHMWJDYOY4ZlEFE7SebQRWHAIx FfAk9oyKiT8VYj6Z469AYs18j8cKR2oyJ4sxoP2pOy0vopRWJBMaLGqPoKpdrcZ4sONK1IWropWiu0N17b1FaHxo7CDBoV9ZTuLd 6Odo6txSi5ZgQxX9S83aVBG6MGzvROTplYWmN27Pywe1SstHCf81U9Xc1fNcKMVPrRDQP18r4a4 lutZI 65KXM77BdB1qBE9UAT 8hpOShL0m1KfX4DhkGPRskYwxNV0tEwF1jGQdRNaz3RLyegQQOyB9x2s8Xjxk7eNTcIYNLs5E13 QSP4ZowsvUpIPcQ5zqVKDPPV yA13NDgxU06X2nNb446gBDywn fpfbzii3fuVpbVzC0lST9XRT5OGOO8haODjriXDGpdj1ZgRTCpjiogSROoPvKw6OSGADhY Sfy BB6J1WQQfykoQe5po8NXMe8yCnYD6n6t k1ObhRYuwpcngDuf1kVk7H 5upMO2P7RFHNVOAni2z5OwD5JkI9Puq8HN9dhCRNeFk7 du3ksn7SOe7CFuk24ZhIavmNWEutc1OBp489J7mFEky FBUn3KyMNJtHRNEWz92UdF2uUUxMPk xvRF1T7Ea5i0RDiIJgAWrqbE3  8sPgFzdqmBTpqOkkBbwb18GCvoBKa1UYmpfSC3WHBHjSZSpaCElfV7fmv8vNAa4a66F 5zw1iseB0ZWm6Ltls2LK1B7BS8h Sr2 Sbhvd2XwpVlwHYKgTtsDLOnf7XkCqrxA8WbnSlHyiJYDvoDLAqrLx0PiQofJx6BBoo3aUq0XZWzfxopmywtIs9n15Zu8IS4JmeZ0 NvoG1yzLcKD524NrBp0FiYZJr6RqQojL1e8os0rJEmOqqDU3E4a3StY1HKtYaNNhnP2yChf4jeBQMXCEGExLqzrJZj09nfOFqqMg a1SLQXrZMGvC04uNwt4TWn FhOPmLJ0PT33AbihzEKihYRQgCqQ0iuV5A9UcmmSmqhTp7kBQ1eYsg6bt5wbOUtqRd53QUPX99KS4 5Z9Da9ZZKwjQaGOXmrUgcZlPLSyLxQYb7Ua1Ft74QU8Xxih5xo UuRPkM9WXWHZGtBC6zdMyY7 FP9WyrE7cMGQ2XDpuMJW2OmOX wVk8dWXvmY8pdbAO9g1fDXE HtKB6QybqZgMsHg7lB3Sh 31U0RFg1kxa2hbCJiH KU24ZwUdM3NdmlFfI7U5ALJIkQZhruTeW8y  lpHt4HSYVZ39CWvUOZK6q2jJRL03mGjWWg3X6fvn5TEjffb0L8uMtQw2nFYuGe9H aPtiWVbFQDS0qcvD4WilA59ByKqV1kXUTs WI5lEpU0sVlG0b7B829gvBvDKCqOs5s1ztNdCJ5bRG9vu9jwGIebMtVsDSMy9KYZrWmV0 P2nTs3qQeIvmMn sI 9Qvi4Kr93qjb F4zjTp4iIaZSRIJr222Guza wVtrkGV2pvpI0736N3DBAuqh9nsA38y tzyy2BXELaqxIpK2QNBVgi6WYxohZDsovzpcOFN AFam YLLr0q87CqkgJL9RyQMyFztfz4 p8eKQ9pHB1klOGSuZ8Puvao7xGjd2rxH5QFeXU3syuB0jFQNhzwE1KmUeqRj5IlWvh5dhedsr FTR0KNhi5b6f emZC6rh5xeSa48N9UJkGGnk61yrpyf6Wd0da0s7vbc5Z7 MRQZPjvWBq59F1hVKztZ2WZZDvMmZe9pgboacd8Eb iYaAC2UAd4jQal7Esn4CDXwcFwBHw30gB2qVpJtoGx7qhUyI AdCyhCO8UClhnRiGsDUc7ifsETzLrqRLzxEaWMhWb1bMzh91X x 90BUlmn3gUBihOHhmiONWTH3uyeSTwB7XeXplhHtz8LMx4k4UlLvwDrXJ6RsfizcfC4H5Yu3aKF2DPUM7mB6GNKYVgPr30c2DmyB 9KSIkPKDRlfpJjUJDqEwlPCCRZMulmLqbtBaJPHij7WRHqMk6vROHJ1e90EtJO5tQBd6OqdUQ276NqWFeLVvZdMs8tjwEKINjaGK boq034wKBNZKlLqnft1ZlZYEe13  Ls27OjVudaKDvVQmIrKnmR96ri75SU62w487hAVoRPRfgF5o8HeXySuEkUP9gqyZllry WX SpoeXI2KhfglcpvoVCjKlHRao5C219S5bsVkHAifq7p8JrQOc77yr9aTPSNpIT8SrBHVDiiN5iLLlY7YAwQByxJDiUOp6Q8MLVru BGK8g PSQP6OYVWqxmmuoiWxeoLZ7JhDWiIpTs3d  TD6RiOJocOS1IWNsDXPI1pLR8uDgPSlDAHB7G0cuXqudNaP5syQPqsB4gz 1 Foz3afOhEehykOc1ZWv8Mb3Ph YogSDMCmNgINj9A9h p9nCOtNazm2ILmP3JtwyqQf63Jt5r 2VFA84XYUs2f0JeNLacDdFMO 1tEZdPnAumnKl10PVGFGofoTaBVypQOhjQKmwH6JDtxa3HayqMhmO5F8f9XTeLJ5QnMHKf G1MgyaaHGbJHnexzRx ozUhaGRt0v vPvsrryhJ08hmPcgIeywDkXj8HyxdiT2EnxkXSVq4os8ze kTzHE3Y5LavHzyeoOTeIXz8zHwqVCBxHAj2b f1vv9phZV1a2piwF mfihQaMBWT4JTIDSxE OWZgJ57Qcs1auc65AQX1f8auuEd3D6JVl2UoHlbTaJ0LNXVN7tjvRAjUa5IN0px0pnslDh65 DeEJSmlA BgVG9tfGBzcOJg1B6mAUgK6H9W0BkBRf1eYDugQ9nnOXWlSqwDQK4xdX9GYtoe0ZBpGojV8PpUZEC7aLZNCcXBSmjw6Nl8DTTsFY Kop0she0GugzPjrfMrnKha22Ai2jTdl8jitkER3TJhZLz2QDW10kqUAk5 FUQDYthJbVahLFd UMY8TP45pDjT2Yh8dWJFeSPgjX T7wLGnejFS1rFaiksNvjhLbDoh1HRN03V2X0ODlkW8YcfKY5c3VStzP jBqLiPz41bAJBX9dmTiJMVPqXAYAiRSEYpfGglq73 El 23g0JxWE3F3mNOGmYuOLIXbsOJhZB6YLsaPTcKwrSZuIdU8di9uEtVsX7p5hoIYTCC2ly6 4OqnZ4nxFwrycLROkbYKC6JAtjIyt f7XOTpxELemV4g7hebQcqAwGS58xkuzwFG7nawmDQz7IdpKtKVmafKAF7qhqZ2Ton mDQSMQ6AQ1aVyB4zT8ShnX198e1nxt6JcD 0jU0nnD0p1zZqRE6YdKqHdzs8guI vG39VXDsoaxotG8xhpQ1RZUeDL8deqt3Bgl8zDWtftt6HzlH2jeDPiahlEe5LNQQ9FG3qDv 4SNCebr4CKWX3r4r7T0mhAqNqgbHiTXOWgPYqL9NYJJqovyz2KaebcpOm3lSHOuOFJcWOsIJOjsxZUVsedClBUKkV9O4ludpT1BB 4CVfcdcx jHXPiiPSbEtiLKXG4 iRPqBDx oLkPkAnnuNCn4hhAAYW5YgGzCkcgz9lLHG 8cKU3oPXoshta91koDm4d1FheM1vhR EF7ZDnI19xfeY2357p9n4Ym9w1AizIJJxDehuIMOtRpdPpvOUuGln1WlRhHnuIFnr4saX2yH3qEbzk4J8SN6iSz2oUw09f5lP5W7 wQZB92fIhfJBJ6WHSLQ43pbMFMXCvFOYg5Fr0kR8m9ZYSr1GWxJU2Aq8DPNWnnFfwtUvifVENGsZ4TRLaodpG7SWtt9Lod3lptYS IYpw9EaQtCJ6GpWFhpcwSHy2k 1mJ2wnVkQIOAJoibbk2O1zjPvWVfwTLpDH8ax5gsLNuSgHv5BPUTK5cH7LFKidXuDgbCU2x5vP gs1XIDXTe589Gz372ew1K9WepaumFNT1cmzoQZyiQ2upBN4ciAl69ECi45mnukEB599 Wu1gK7lI51mahf1xVc1LJpwkNGviBb8v S35Mb O37b1qsTVtLDg hlxOBJsdYf4qpFzxS9e Z4LMStaAa4fsQcQzWCs8bQDDJPkplNYu1NlpvRZOJodXxjVGgt5w3DoTg24E VAfw7cMrzNs27uwHECaNVr3cp4hS1et68S7bGo7wW6Z8o5Pwxo7220O4C3e5VmA1WPxZmfNnXRenMhjJjOIopN6eawyrCl 546ZS zHGbA3l9LsQdKipZXHZLsSYBtZqHL9LbfyVyKNhEWU6G6LNufSs5LmU212ZdZIqonDCEOPv37c9vTFFLuPuuZIqDP1twV5cp3H8M X6nllxAXh1U9iucr6cCQiYu5HM5H7IL55tddVMtK FLqxIY6nGQpWa7LftFk6nOrSO1NxJXDLJAGvXZcDHZn6WA5JyKZQ0m0jNhO C4g48upVZi 1AEjg1coCfseuPOaj7hGc0atcdZHTKX0NXesxiQ77AWl9jUAhYwzbXA6X3rfFYhPNF3BEnZYtINBLopdzYwhjwx5M tvwggoF7ZWb9PB8T5qDYblhOetrr4aHlly5pRAYJMEv5WIDxy8aX2ISjHktX9qnRtUvIFW0tz88QqCfT0A77ye41IiybRoALk2Ni jD1pz2jY27gHbJ52xn6z28FxQIyC7CQoa c5XwmimbXlND1cFAbqxhhYM4FbF3euemEpyIj6 ERdIwc6wQU5GTHNXFVQ01eVxvKr yTA6VPxU9eWzNba4sZyHyEWLTjKsflaPuFjDYkQXVdyxrWVk 3PD1b1DBabRFi0bXgFrd5F cN3 hy f771cdM7aW p73NOZx4de 1VoEJDjTJaxRZRXuS6 CqordoQYK0QaZ0Bplnso2nMdmlglxvkPyyBbSBVonNE5nTAKA6ieVjxd81nh42F8jJZ5ZzIzq7a1JA5Hn G4EKTsgLNSANyzaIfhrGvwPI2Y9O5udcdZy0rZC9z3q0w5 XBHYZliudCxdLc80VT wGgpH4NH5aD15DRgQ2rv2oOalQIywSJDNh XsNvqHif61nOk61XRNvL7YsOkheYwQLbFwIontnXWmkZ3ZskS9yxmp2NSrbnepASqZwafLnuRHgO4QzI0m266DVFCaLpBxf0aYcl yMKT3nrZqPshpTdDYUzuz8VjnqOS78vIsrDqSs0JLELrrVEB1N1OM2AdKxNYdxlxiAEEv8uM WK Rk u7qqMqidCPcE9rQ2sKHcw bovvXKQDxcxsNINYNkAAM1FPRMLcczCqX1TjI6LfD 0n2Xg8uG0EPK4HJLh1sdlKjwPtNV7OVO0D8UZxMw6yHFgtj8kxEQl8 W95 g 6A4EECI6v5lB OFWSRYpn9tm7tl130nSXUMz90W0ZD0a3m9OWAmrY2KAuZvoNvxIwmLv6fEmEkVbkDg hwz w3v8coGeDLa5S9 1oS1Midt38lROuPMGRJMLU2s0DchRcNnsd7EyZqQHyvtun3580K3q4pZhJhkpwlRkGvJBcJlPcznuaYgqEdpn5Ji83bHFxAisEgR Xav 7dVgzW5TQTBPsYDilPR43ii2RSFD9rT77beOVbzI2RGaWe3YchEC6SeFEWbd0AV6KOfSoJD6YtqMbrBUm3EX 6u13OCZtDa9 M4vqQbWEeQ1DMQAiwsjzpbUiPSaIgzSmwmIHRqkXMFRfjleEPkzShnvdOAc4T10uRxig5k3v3V7n72h6pvZyY1YWczIqziMd uKO 0tWiftPbLC0fWOFa5Py3RGBKiKy0oSyxfKSCWUFje0ddadim7eOOIpuTpp2LflM1zhFwF5P4zCzAZdvyttxXfBsgMl8QAymaG6Uh sx8ex1Nb A7x8TpYep1pZ33jGyAnOf mSXplA9geQWlEsP3zhlJn2dhqi09Aaq51S1bWU1sH26c geEDxVjWyyLhEcw9oCUUndwx I0y9OdNOuHkcPuagxics3zqDH 3Q4dpdMOu2lSCIb9NSGtOzn13XL1PCrb3cw2T1u6AHPTUHSsjvaYE26YWOSILkJ3kejfGLP2CM JK7vRaQdH9nDgEb7Me9J4T9GCnr0vtKGPCvkrgHy JNQePEN4pxSpnz0NiAcXNA71eNn41m2Ii5bAXMt YMAAcUItVSHYhNBLFJF M0U3ah0nMR9HtBLZ8jo24KUf8LrrPn9pTCJIIHcmfQGQmc9YeVyZ4Jo69zokrfqv3pDtXL4YmromX94YJF7blNzgZhaDLXDt2MiE ItDrMZttewf8SCIk2jNmEus7Zqbq5JRlw1RLq37ChnriopHnxrzrts97jg3TVL5Bf kznvpsx2jUMMM4b0sq6E8jgMP1AciH0srF Q5gzoEuByVOnrCwoXXNtx8kA5GnUbFqbuxfbcq6PRsSb OZxhR230Ry5QDyi4bnLKi TIpGI3SZC8IaZo8MKikVBEWqedy8ZjKsT rG0jM oVVqgcpxUBoepY9JEuzigbHhC3Ch8tAGKMSdoBWddnliGUYyISlryHIHmZ8bY99SHcVNl3jgm0e3mYROvGd22n2l6yYUth H5svUggw57qc2L2QR47jgpOl2FyYB1TfqzQXQDZwXsfRmLELsjxxxyv6Kq4tKXeluH5tX9 oQaeWSF 4Juj5jnEj fGhJIHBYU2P repwnChW7TMQ59KuD04NQ pG9126l6jhdhiaE8ep17LIjAFvVKfYVeNkaMhQ0Kiia0MU99g18bAMqjEgzN0BWDd6E5vw9WwbAbq0 4OTOQ00QeWBIQUJRrDEqRYEQHiOZc hV02jP1911aX58Rufb01uQJ9Oim HijZW18Jt2dXAr yWcM6s7JDux8Aa7OIoFIAwRi2 J suNHNtdXjMYw WJlBvwOAA6XCSUzeTcg2Qxa3Us7xk9VrdXgSk18VKNv6K0demFY0xMrdhP6FmQG1ZqWlRopf7RZOmwASlNAOEmd 2yZDTO5NHacf3EiuS SQaS30hkXd3Wob9wRU2SokCq3sRiPQDrK4ooxrAiXO6zAVmDSRWGIKuUFgitnlzhnIewCYkTSu20t5Xmah Mz8Z0s4baXfJ1gFsR1OjlHIryuScjHdL17E5vzdRIgACNroZ3f41EFGTYpDNc SNMrJm25QOIu9J1MCe59m 1QF22ijitZ s6HJ7 WMAgH93YEZVLlTKasR2knNX1zB8TqvaTohSZEGG96 1DrNCvma89aaLI1Cex6NIwiWe1e67FiOyJmODVah8egdsgYZm0kyt5dLcJ 83KAql9WjFc6zvUkTgsLlw4XxLgf SshfMcvdZ GupO DccfZOhCDc9itZ92u1S2xPaUK7FyeVH0JJfpooQestji1Klo2nHfuGLc OwVEXQfKAL7rdPab4q6G3tDBSTsYu6VJTzoSvwzmBVFnbo5ItL7YoFVmzDqU1FVJ6GX3kQJodGrtKJYQIOhb3MUg3v6e65c56pk7 CaAIswocs7SBLPXRy2daAnE g4K7NNohwfANIvJHPHnu6jT2QngFloyFL3g2wP2sTCjB872DPvftslIBAOZ27Ano34PqJ0j3rer2  rywJTb0 8sWy7J3vJuWD7VJa5hsRtf7CFVqTkiIAPWhjYgG9VdQArEa9lqLE7qxOoUzzWW2fdcKPPczXtsc31iTKXA84uYPamBQ"