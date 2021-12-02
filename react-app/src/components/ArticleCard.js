import React from "react";
import ArticleCardActions from "./ArticleCardActions";

const useStyles = makeStyles((theme) => ({
    editButton: {
        marginRight: 8,
        display: 'inline',
    },
    cancelButton: {
        float: 'right',
        color: '#34495E',
    },
    gridContainer: {
        paddingTop: theme.spacing(4),
    },
    importantCell: {
        fontWeight: 600,
        fontSize: 18,
    },
    infoIcon: {
        verticalAlign: 'middle',
    }
}));

export default function ArticleCard({ document, removeCitation, addCitation }) {
    const theme = useTheme();
    const classes = useStyles(theme);

    const citationRow = (document) => {
        return (
            <TableRow className="article-card" key={document.id}>
                {/* Title */}
                <TableCell className={classes.importantCell}>{document.Title}</TableCell>

                {/* Author(s) */}
                <TableCell>
                    {document.Authors.map((author) => author.Name).join(", ")}
                </TableCell>
  
                {/* Tooltip Blurb? */}
                {/* <TableCell>{row.WrapRateMultiplier ? row.WrapRateMultiplier : defaultWrapRateMultiplier}</TableCell> */}

                <TableCell align="right">
                    <ArticleCardActions
                        document={document} 
                        addCitation={addCitation} 
                        removeCitation={removeCitation} />
                </TableCell>
            </TableRow>
        );
    }

    return (
        <React.Fragment>
            {
                // citationTable.length > 0 ?
                (<Table size="medium">
                    <TableHead>
                        <TableRow>
                            <TableCell width="16.25%">Title</TableCell>
                            <TableCell width="12%">Author(s)</TableCell>
                            <TableCell width="7.5%">Blurb</TableCell>
                        </TableRow>
                    </TableHead>
                </Table>) 
            }
        </React.Fragment>
    );
}
