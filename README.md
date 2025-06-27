# City Population Service

This is a simple HTTP API written in Go that allows users to store and retrieve population data for cities using [Elasticsearch](https://www.elastic.co/elasticsearch/).

## Features

- **Health Check**: `/health` endpoint for basic service availability.
- **Add/Update City**: `/city` endpoint to add or update city population data.
- **Query Population**: `/population?name={city}` to retrieve population data for a given city.
- **Data Storage**: Uses Elasticsearch for persistent storage.

---

## Requirements

- Go 1.18+
- Elasticsearch 7.x or 8.x running locally or remotely

---

## Environment Variables

| Variable       | Description                             | Default                  |
|----------------|-----------------------------------------|--------------------------|
| `ES_URL`       | Elasticsearch URL                       | `http://localhost:9200` |
| `ES_USERNAME`  | Username for basic auth (if required)   | `""` (empty)             |
| `ES_PASSWORD`  | Password for basic auth (if required)   | `""` (empty)             |

---

## API Endpoints

### `GET /health`

Returns a simple `OK` response to confirm the server is running.

**Response:**
```
OK
```

---

### `POST /city`

Stores or updates a city's population data in Elasticsearch.

**Request Body:**
```json
{
  "name": "Berlin",
  "population": 3769000
}
```

**Response:**
```
City stored/updated successfully
```

---

### `GET /population?name={city}`

Fetches the population data for the specified city.

**Example Request:**
```
GET /population?name=Berlin
```

**Example Response:**
```json
{
  "name": "Berlin",
  "population": 3769000
}
```

---

## Running the Server

1. Set environment variables (if needed).
2. Start Elasticsearch locally or ensure remote access is available.
3. Run the server:

```bash
go run main.go
```

Server starts on port `8080`.

---

## Notes

- City documents are indexed in the `cities` index.
- Document IDs are lowercase city names to ensure consistency.

---


## Deploying via Helm Chart

The helm chart is packaged as an OCI compatible image and can be found here https://github.com/jeicoo/population/pkgs/container/population-helm-charts%2Fpopulation

### Prerequisites

Before deploying the helm chart, you must have access to a cluster which meets the following criteria:

- Kubernetes version 1.31 and up
- CSI Driver configured for Persistence
- Elastic Cloud on Kubernetes(ECK) operator installed in the cluster. [Read more](https://www.elastic.co/docs/deploy-manage/deploy/cloud-on-k8s)

### Deploying

1. Download the chart
    ```bash
    helm registry login ghcr.io
    helm pull oci://ghcr.io/jeicoo/population-helm-charts/population --version v0.0.10
    ```

2. Verify the chart
    ```bash
    helm show all oci://ghcr.io/jeicoo/population-helm-charts/population --version v0.0.10
    ```
