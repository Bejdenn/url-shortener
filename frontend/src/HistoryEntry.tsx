import { Link, Stack, Typography } from '@mui/material';
import { URLRelation } from './types';

const HistoryEntry = ({ relation }: { relation: URLRelation }) => {
  return (
    <Stack margin={2} spacing={2} direction='row' alignItems='center' justifyContent='space-between'>
      <Typography>{relation.longUrl}</Typography>
      <Link sx={{ fontFamily: 'sans-serif' }} rel='noopener noreferrer' target='_blank'
            href={relation.shortUrl}>{relation.shortUrl}</Link>
    </Stack>
  );
};

export default HistoryEntry;