# bt-whole-home-exporter

This project exports devices, signal strength and system metrics from the BT Whole Home Wifi system. They're put on a `/metrics` HTTP endpoint that you can scrape with Prometheus.

You can run it using:

```sh
go build .
./bt-whole-home-exporter $YOUR_BT_WHOLE_HOME_WIFI_ADMIN_PANEL_PASSWORD_HERE
```

It assumes your BT Whole Home WiFi's admin panel is running on `192.168.1.1`.

## Maintainance warning

This project is not maintained. I used it for several months and it was extremely reliable with the `Whole Home Wi-Fi v1.02.12 build02` firmware version.

The code is a bit terrible because (a) some strange decisions in the admin panel's security (b) it was built in a hurry.

Also worth noting is that (at the time of writing) there's an [active security issue affecting the latest firmware](https://community.bt.com/t5/BT-Devices/BT-Wi-Fi-Disc-susceptible-to-authentication-bypass/td-p/2177297).
