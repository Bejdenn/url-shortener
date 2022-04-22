import LoadingButton from '@mui/lab/LoadingButton';
import { Input, Stack, SxProps, Theme, Typography } from '@mui/material';
import { useState } from 'react';
import './App.css';
import History from './History';
import { URLRelation } from './types';
import useLocalStorage from './hooks/useLocalStorage';
import 'react-toastify/dist/ReactToastify.min.css';
import { toast, ToastContainer } from 'react-toastify';
import Layout from './Layout';

export let API_URL: string;
if (process.env.NODE_ENV === 'production') {
  API_URL = 'https://api-72ey6bex.nw.gateway.dev';
} else {
  API_URL = 'http://127.0.0.1:8080';
}

function App () {
  const [isLoading, setIsLoading] = useState(false);
  const [urlToSubmit, setUrlToSubmit] = useState('');
  const [history, setHistory] = useLocalStorage<URLRelation[]>('urls', []);

  const getShortUrl = async () => {
    setIsLoading(true);

    const response = await fetch(API_URL + '/url-shortening', {
      method: 'POST',
      body: new URLSearchParams({ longUrl: urlToSubmit })
    });

    setIsLoading(false);

    if (response.status === 200) {
      const relation = await response.json() as URLRelation;
      setHistory([relation, ...history]);
      setUrlToSubmit('');

    } else if (response.status === 400) {
      toast.error('Your URL is not valid.', {
        position: 'top-center',
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        progress: undefined
      });
      throw new Error(await response.text());

    } else {
      toast.error('An error has occurred. Please try again.', {
        position: 'top-center',
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        progress: undefined
      });
      throw new Error(await response.text());
    }
  };

  const fontVariant = 'h4';
  const fontStyle: SxProps<Theme> = { fontWeight: 'bold', color: 'white' };

  return (
    <Layout>
      <Stack spacing={4}>
        <Typography sx={{ ...fontStyle }} variant={fontVariant}>
          Shorten your link 🔗
        </Typography>
        <Input
          value={urlToSubmit}
          sx={{ color: 'white' }}
          placeholder="URL address"
          onChange={(e) => setUrlToSubmit(e.target.value)}
        />
        <LoadingButton
          disabled={urlToSubmit.length === 0}
          variant="contained"
          loading={isLoading}
          onClick={getShortUrl}
        >
          Shorten
        </LoadingButton>
        <History isDev={process.env.NODE_ENV !== 'production'} entries={history}/>
      </Stack>
      <ToastContainer/>
    </Layout>);
}

export default App;
