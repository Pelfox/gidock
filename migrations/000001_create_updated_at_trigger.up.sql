-- 1. Create the regular trigger function
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Create the event trigger function that attaches the trigger to new tables
CREATE OR REPLACE FUNCTION auto_add_updated_at_trigger()
    RETURNS event_trigger AS $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT * FROM pg_event_trigger_ddl_commands()
        LOOP
            -- Only process CREATE TABLE commands
            IF r.command_tag = 'CREATE TABLE' AND r.object_type = 'table' THEN
                -- Check if the new table has an updated_at column of type TIMESTAMPTZ
                IF EXISTS (
                    SELECT 1
                    FROM information_schema.columns
                    WHERE table_schema = COALESCE(r.schema_name, 'public')  -- fallback to 'public' if no schema
                      AND table_name = split_part(r.object_identity, '.', array_upper(string_to_array(r.object_identity, '.'), 1))
                      AND column_name = 'updated_at'
                      AND udt_name = 'timestamptz'
                ) THEN
                    EXECUTE format(
                            'CREATE TRIGGER set_updated_at
                             BEFORE UPDATE ON %s
                             FOR EACH ROW
                             EXECUTE FUNCTION trigger_set_updated_at()',
                            r.object_identity
                    );
                END IF;
            END IF;
        END LOOP;
END;
$$ LANGUAGE plpgsql;

-- 3. Create the event trigger itself
CREATE EVENT TRIGGER auto_updated_at_trigger
    ON ddl_command_end
    WHEN TAG IN ('CREATE TABLE')
EXECUTE FUNCTION auto_add_updated_at_trigger();