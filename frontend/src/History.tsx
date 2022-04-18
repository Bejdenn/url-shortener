import { Divider, Stack } from '@mui/material';
import { URLRelation } from './types';
import HistoryEntry from './HistoryEntry';

interface HistoryProps {
  entries: Array<URLRelation>;
}

const History = ({ entries }: HistoryProps) => {
  return (
    entries.length > 0 ?
      <Stack sx={{ bgcolor: 'white', borderRadius: 1 }}
             divider={<Divider orientation='horizontal' flexItem />}>{
        entries.map(rel => <HistoryEntry key={rel.id} relation={rel} />)
      }</Stack> : <></>
  );
};

export default History;