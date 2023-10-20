BEGIN;
  CREATE TABLE IF NOT EXISTS metrics(
     id VARCHAR(50) PRIMARY KEY,
     mtype VARCHAR(7) NOT NULL,
     delta BIGINT NULL,
     value DOUBLE PRECISION NULL,
     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
  );

  COMMENT ON TABLE metrics IS 'metrics storage';
  
  COMMENT ON COLUMN metrics.id IS 'Unique metrics ID';
  COMMENT ON COLUMN metrics.mtype IS 'Metrics type gauge or counter';
  COMMENT ON COLUMN metrics.delta IS 'Counter type metrics value';
  COMMENT ON COLUMN metrics.value IS 'Gauge type metrics value';
  COMMENT ON COLUMN metrics.created_at IS 'Row created date';
  COMMENT ON COLUMN metrics.updated_at IS 'Row updated date';
COMMIT;
