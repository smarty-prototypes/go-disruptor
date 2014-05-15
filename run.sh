#!/system/bin/sh

chmod 755 ./go-disruptor
while ./go-disruptor 2>/dev/null; do
	# sleep 1
done
