# Nick Selective Sync Sensor

### Build:
```
make
```

### Example Config:
```json
{
  "components": [
    {
      "name": "selective",
      "namespace": "rdk",
      "type": "sensor",
      "model": "ncs:sensor:nick-selective-sync-sensor",
      "attributes": {"sync": true}
    }
  ],
  "services": [
    {
      "name": "data_manager-1",
      "namespace": "rdk",
      "type": "data_manager",
      "attributes": {
        "additional_sync_paths": [],
        "sync_disabled": true,
        "maximum_num_sync_threads": 10,
        "sync_interval_mins": 0.01,
        "capture_dir": "",
        "selective_syncer_name": "selective"
      }
    }
  ],
  "modules": [
    {
      "type": "local",
      "name": "selectivesensor",
      "executable_path": "/home/user/nick_selective_sync_sensor"
    }
  ]
}
```
