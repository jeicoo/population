> For helm chart deployment, go to [Deploying via Helm Chart](#deploying-via-helm-chart)
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

The Population Helm Chart is packaged as an OCI-compatible image and is available at:  
🔗 [https://github.com/jeicoo/population/pkgs/container/population-helm-charts%2Fpopulation](https://github.com/jeicoo/population/pkgs/container/population-helm-charts%2Fpopulation)

---

### Prerequisites

Before deploying, ensure the following are in place:

- **Helm 3**
- **Kubernetes** v1.31 or later
- **CSI Driver** configured for persistent storage
- **Elastic Cloud on Kubernetes (ECK)** operator installed  
  📘 [ECK documentation](https://www.elastic.co/docs/deploy-manage/deploy/cloud-on-k8s)

---

### Deployment

To install or upgrade the Population Helm release:

```bash
helm upgrade --install population oci://ghcr.io/jeicoo/population-helm-charts/population --version 0.0.11
```

By default, this installs:

- A **single-node Elasticsearch** instance with **self-signed TLS certificates**
- The **Population services**, preconfigured to connect to the bundled Elasticsearch

The chart includes the [`eck-elasticsearch`](https://github.com/elastic/cloud-on-k8s/tree/main/deploy/eck-stack/charts/eck-elasticsearch) subchart. For customization options, refer to its `values.yaml`.

The application requires the following environment variables to connect to Elasticsearch:

- `ES_URL`
- `ES_USERNAME`
- `ES_PASSWORD`

You **do not** need to set these manually—the Helm chart injects default values automatically.

---

### Production Deployment

> ⚠️ **Recommendation:** For production environments, use a separately managed Elasticsearch cluster.

While the bundled single-node Elasticsearch is suitable for testing or development, it's not ideal for production. Storage services like Elasticsearch typically have separate lifecycles and scaling requirements, which are better handled independently from the application.

To disable the built-in Elasticsearch:

```yaml
elasticsearch:
  enabled: false
```

---

### Configuring an External Elasticsearch Cluster

When using an existing Elasticsearch deployment, provide the connection details via the Helm values:

```yaml
# Used only if elasticsearch.enabled is set to false
config:
  elasticsearchUrl: "https://elasticsearch:9200"
  elasticsearchUsername: "elastic"
  elasticsearchPassword: "changeme"

  existingSecret:
    enabled: false               # Set to true to use an existing Kubernetes Secret
    secretName: "secretname"     # Replace with the name of your secret
    passwordKey: "elastic"       # Key within the secret containing the password
```

---

### Additional Configuration

For more advanced settings and customization options, refer to the chart’s `values.yaml`:  
📄 [values.yaml](https://github.com/jeicoo/population/blob/main/chart/values.yaml)
