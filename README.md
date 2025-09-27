# Weird prometheus exporters

This repository contains some weird prometheus exporters that I have built either because they are idiosyncratically useful to me or I find them funny (or both).

They are not:

* good code
* finished
* embodiments of best practices
* good things to infer my competence (or lack thereof) from

Don't use them for anything important.

* `appletalk_exporter`: uses netatalk 2's utilities to export information about what entities are on an appletalk network.

* `backupdir_watcher`: given a directory of files with dates embedded in their names, is there a recent one?

* `dir_count_watcher`: how many files are in a directory?

* `env_agency_flood_exporter`: make available the height of rivers using environment agency open data

* `findmy_battery_exporer`: this one is some next level bullshit, runs on a macOS virtual machine and makes available battery levels of iCloud connected devices.

* `lms_exporter`: Stats about Lyrion/Logitech Media Server

* `realtime_trains_exporter`: exports stats about trains through a station.  You need a realtimetrains API account.

* `sas2ircu_exporter`: SAS-2 integrated RAID utility output

* `tinydns_exporter`: more than you wanted to know about tinydns

* `truenas_api_exporter`: a couple of things that are missing from the collectd stats: whether there are any active alerts, and how cloud sync jobs are doing.

* `usb_device_exporter`: lsusb, but for prometheus

* `wiltshire_bins_exporter`: do you live in wiltshire and want to know whether you should put your bins out tonight?

* `xset_monitor_status_exporter`: has your monitor turned off courtesy of DPMS?

* `zigbee2mqtt_metadata_exporter`: what Zigbee devices are configured and how they are connected to one another.
