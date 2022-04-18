import LoadingButton from '@mui/lab/LoadingButton';
import { Box, Input, Stack, SxProps, Theme, Typography } from '@mui/material';
import { useState } from 'react';
import './App.css';
import History from './History';
import { URLRelation } from './types';

const API_URL = 'https://api-72ey6bex.nw.gateway.dev';

const App = () => {
  const [urlToSubmit, setUrlToSubmit] = useState('');
  const [history, setHistory] = useState<URLRelation[]>([]);

  const getShortUrl = async () => {
    const payload = { longUrl: urlToSubmit };
    const response = await fetch(API_URL + '/url-shortening', { method: 'POST', body: new URLSearchParams(payload) });

    if (response.status !== 200) {
      throw new Error('error in response: ' + await response.text());
    }

    const relation = await response.json() as URLRelation;
    setHistory([relation, ...history]);
    setUrlToSubmit('');
  };

  const fontVariant = 'h4';
  const fontStyle: SxProps<Theme> = { fontWeight: 'bold', color: 'white' };

  return (
    <Box sx={{ paddingTop: 4, paddingLeft: 2, paddingRight: 2 }}>
      <Box display={'flex'} justifyContent='center' alignItems='center'>
        <Box
          sx={{
            paddingTop: '6rem',
            maxWidth: '50%'
          }}
        >
          <Stack spacing={4}>
            <Typography sx={{ ...fontStyle }} variant={fontVariant}>
              Shorten your link 🔗
            </Typography>
            <Input
              value={urlToSubmit}
              sx={{ color: 'white' }}
              placeholder='URL address'
              onChange={(e) => setUrlToSubmit(e.target.value)}
            />
            <LoadingButton
              disabled={urlToSubmit.length === 0}
              variant='contained'
              onClick={() => getShortUrl()}
            >
              Shorten
            </LoadingButton>
            <History entries={history} />
          </Stack>
        </Box>
      </Box>
    </Box>
  );
};

export default App;
