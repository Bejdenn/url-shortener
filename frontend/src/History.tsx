import { Divider, Stack } from '@mui/material';
import { URLRelation } from './types';
import HistoryEntry from './HistoryEntry';

interface HistoryProps {
  entries: Array<URLRelation>;
  isDev: boolean;
}

const History = ({ entries, isDev }: HistoryProps) => {
  return (
    entries.length > 0 ?
      <Stack sx={{ bgcolor: 'white', borderRadius: 1 }}
             divider={<Divider orientation='horizontal' flexItem />}>{
        entries.map(rel => <HistoryEntry isDev={isDev} key={rel.id} relation={rel} />)
      }</Stack> : <></>
  );
};

export default History;