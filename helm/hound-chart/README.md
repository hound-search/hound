# Helm chart for 'hound'

## Configuration

The configuration of the repositories to scrape is mounted in a secret due it can contain credentials.

```
helm install ./helm/hound-chart --set efsFileSystemId=fs-12345678 --set awsRegion=ue-central-1
```

