# SmartLoad Optimization API

Core microservice for matching trucks with the most profitable combination of shipments.

## Key Features
- **Optimal Solver**: Recursive backtracking with pruning for $n=22$ orders.
- **Constraints**: Handles multidimensional capacity (weight/volume) and Hazmat isolation.
- **Production Ready**: Multi-stage Docker build, structured logging, and health checks.

## How to Run
```bash
git clone [https://github.com/saptaka-trihantoro/optimal-truck-load-planner.git](https://github.com/saptaka-trihantoro/optimal-truck-load-planner.git)
cd optimal-truck-load-planner
docker compose up --build
```
## Testing
I have included both Unit Tests for the algorithm and E2E tests for the API layer.
To run all tests:
```bash
go test ./... -v

## API EndpointsPOST ## 
/api/v1/load-optimizer/optimize : Core optimization engine.GET /healthz : Service health status.

## PerformanceTime ##
Complexity: $O(2^n)$ with heavy pruning.Average latency for $n=22$: < 50ms.