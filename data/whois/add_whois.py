#!/bin/env python3
# Whois Server: https://github.com/spdir/whois-server-list

import json
import socket
import threading

trs = []

ds = open("whois.servers.json", 'r')
ds_json = json.loads(ds.read())
ds.close()


def connetcTest(url, k):
	try:
		sk = socket.socket()
		sk.connect((url, 43))
		send_str = "test" + "." + "k"
		sk.send(send_str.encode("utf-8"))
		sk.send("\r\n".encode("utf-8"))
		recv = sk.recv(10240)
		# print(recv)
		print(url, "Yes")
		sk.close()
		return True
	except Exception as e:
		# print("error", e)
		print(url, "No", " error:", e)
		return False


def test_whois_server(k, v):
	if not ds_json.get(k, None):
		ds_json[k] = []
	if connetcTest(v, k):
		whois_list = ds_json[k]
		whois_list.append(v)
		whois_list = list(set(whois_list))
		ds_json[k] = whois_list


def optionSelfFile():
	"""
	自定义新的文件处理，添加到whois.servers.json
	"""
	f = open("twhois.servers.json", "r")
	whois_server_file = json.loads(f.read())
	for k, v in whois_server_file.items():
		# v = v.get("host", None)
		v = v[0]
		if not v:
			continue
		t = threading.Thread(target=test_whois_server, args=(k, v,))
		trs.append(t)


if __name__ == "__main__":
	optionSelfFile()

	for t in trs:
		print("[%d]-->[%d]" % (len(trs), trs.index(t)))
		t.start()
		while True:
			if (len(threading.enumerate()) < 15):
				break
	dump_data = json.dumps(ds_json, indent=2)
	with open("whois.servers.json", "w") as f:
		f.write(dump_data, )
