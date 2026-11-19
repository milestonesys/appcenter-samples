import React from 'react';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Alert,
  Box,
  Button,
  CssBaseline,
  Stack,
  TextField,
  ThemeProvider,
  Typography,
} from '@mui/material';
import { darkTheme, lightTheme } from '@milestone/milestone-theme';
import { useSharedStore } from '@milestone/cross-mfe-communication';
import type { TMfeProps } from '@milestone/one-shell-mfe-contract';
import { createEvent, deleteEvent, getEvent, getServerInfo, updateEvent } from './api';
import type { ApiResult } from './api';

type Step = { title: string; result: ApiResult };

type JsonPanelProps = { step: Step };

const JsonPanel: React.FC<JsonPanelProps> = ({ step }) => (
  <Accordion>
    <AccordionSummary>
      <Typography variant="bodyLarge">
        {step.title} ({step.result.status})
      </Typography>
    </AccordionSummary>
    <AccordionDetails>
      <Box component="pre" sx={{ margin: 0, overflow: 'auto' }}>
        {JSON.stringify(step.result.data, null, 2)}
      </Box>
    </AccordionDetails>
  </Accordion>
);

// Pull an error message out of the gateway's {"error": ...} responses.
const describeError = (result: ApiResult): string => {
  const data = result.data as { error?: unknown; detail?: unknown } | string | null;
  if (data && typeof data === 'object') {
    const message = data.error ?? data.detail;
    if (message !== undefined) {
      return typeof message === 'string' ? message : JSON.stringify(message);
    }
  }
  return `HTTP ${result.status}`;
};

const App: React.FC<TMfeProps<{}>> = () => {
  const themeMode = useSharedStore((state: any) => state.themeMode);
  const theme = React.useMemo(() => (themeMode === 'dark' ? darkTheme : lightTheme), [themeMode]);

  const [serverUrl, setServerUrl] = React.useState<string>();
  const [eventName, setEventName] = React.useState('');
  const [newEventName, setNewEventName] = React.useState('');
  const [eventId, setEventId] = React.useState<string>();
  const [steps, setSteps] = React.useState<Step[]>([]);
  const [error, setError] = React.useState<string>();
  const [success, setSuccess] = React.useState<string>();
  const [busy, setBusy] = React.useState(false);

  React.useEffect(() => {
    getServerInfo()
      .then((result) => result.ok && setServerUrl(result.data.serverUrl))
      .catch(() => undefined);
  }, []);

  // Runs one backend call, turning thrown network errors into the error banner.
  const run = async (action: () => Promise<void>) => {
    setBusy(true);
    setError(undefined);
    setSuccess(undefined);
    try {
      await action();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  };

  const runDemo = () =>
    run(async () => {
      setSteps([]);

      const created = await createEvent(eventName);
      setSteps([{ title: 'Create item result', result: created }]);
      if (!created.ok) {
        setError(`Failed to create event: ${describeError(created)}`);
        return;
      }
      const id = created.data.result.id;
      setEventId(id);

      const fetched = await getEvent(id);
      setSteps((previous) => [...previous, { title: 'Get item result', result: fetched }]);
      if (!fetched.ok) {
        setError(`Failed to get event: ${describeError(fetched)}`);
      }
    });

  const runUpdate = () =>
    run(async () => {
      if (!eventId) return;
      const updated = await updateEvent(eventId, newEventName);
      setSteps((previous) => [...previous, { title: 'Update item data', result: updated }]);
      if (updated.ok) {
        setSuccess('Event updated successfully!');
      } else {
        setError(`Failed to update event: ${describeError(updated)}`);
      }
    });

  const runDelete = () =>
    run(async () => {
      if (!eventId) return;
      const deleted = await deleteEvent(eventId);
      if (deleted.ok) {
        setSteps([{ title: 'Delete response', result: deleted }]);
        setEventId(undefined);
        setSuccess('Event deleted successfully!');
      } else {
        setError(`Failed to delete event: ${describeError(deleted)}`);
      }
    });

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Stack spacing={2} sx={{ padding: 3, maxWidth: 800 }}>
        <Typography variant="headingLarge" component="h1">
          RESTful Config API
        </Typography>
        <Typography variant="bodyMedium" color="text.secondary">
          Calls to the XProtect API Gateway are made on behalf of the user signed in to this Smart Client session.
        </Typography>
        {serverUrl && <Alert severity="info">XProtect server URL: {serverUrl}</Alert>}

        <TextField
          label="Name of user defined event to create"
          value={eventName}
          onChange={(e) => setEventName(e.target.value)}
          fullWidth
        />
        <Box>
          <Button variant="contained" onClick={runDemo} disabled={busy || !eventName}>
            Create Event
          </Button>
        </Box>

        {error && <Alert severity="error">{error}</Alert>}
        {success && <Alert severity="success">{success}</Alert>}
        {steps.map((step, index) => (
          <JsonPanel key={`${step.title}-${index}`} step={step} />
        ))}

        {eventId && (
          <Stack spacing={2}>
            <TextField
              label="Enter new event name"
              value={newEventName}
              onChange={(e) => setNewEventName(e.target.value)}
              fullWidth
            />
            <Stack direction="row" spacing={2}>
              <Button variant="outlined" onClick={runUpdate} disabled={busy || !newEventName}>
                Update Event
              </Button>
              <Button variant="outlined" color="error" onClick={runDelete} disabled={busy}>
                Delete Event
              </Button>
            </Stack>
          </Stack>
        )}
      </Stack>
    </ThemeProvider>
  );
};

export default App;
