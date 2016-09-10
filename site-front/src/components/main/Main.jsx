import React, { Component } from 'react';
import { connect } from 'react-redux';

import LanguageSelector from './LanguageSelector.jsx';
import FixHeaderBox from './FixHeaderBox.jsx';

import ConnectFourBoard from './ConnectFourBoard.jsx';
import './Main.scss';

import { newGame, insertToken } from '../../actions/ConnectFourActionCreators.js';

const langs = ['german', 'english']

export class Main extends Component {
    
    constructor(props) {
        super(props);
        this.state = { 
            content: undefined,
            hideFixHint: false
        };
        this.switchLanguage = this.switchLanguage.bind(this);
        this.fixHeader = this.fixHeader.bind(this);
        this.play = this.play.bind(this);
    }

    componentWillMount() {
        this.switchLanguage('german');
        //this.props.newGame();
    }

    switchLanguage(lang) {
        if (lang == undefined || lang == '') {
            return;
        }
        const use = this.props.content[lang]
        if (use != undefined) {
            this.setState( {'content': use});
        }
    }

    fixHeader() {
        this.setState({ hideFixHint: true });
        this.props.fixHeader();
    }

    play(column) {
        const { gameId, insertToken } = this.props;
        insertToken(gameId, column);
    }

    render() {
        const { content, hideFixHint } = this.state;
        const { board, cols, rows, newGame } = this.props;
        if (content == undefined) {
            return null;
        }
        return (
            <main className='about'>
                <LanguageSelector selectCallback={this.switchLanguage} langs={langs}/>
                <h1>FELIX ORTMANN</h1>
                <div className='about__subtitle'>
                    &mdash; {content.title} &mdash;
                </div>
                <FixHeaderBox fixIt={this.fixHeader} hide={hideFixHint} displayText={content.box} />
                <img className='about__img'/>
                <div className='about__description'>
                    {content.description.map((paragraph, i) => {
                        return <p key={i}>{paragraph}</p>
                    })}
                </div>
                <object className='gopher__left'/>
                <button onClick={newGame}>Ding Dong</button>
            </main>
        );
    }
}

Main.propTypes = {
    content: React.PropTypes.object.isRequired,
    fixHeader: React.PropTypes.func.isRequired
}

function mapStateToProps(state, ownProps) {
    return {
        gameId: state.gameId,
        board: state.board,
        cols: state.columns,
        rows: state.rows
    }
}


export default connect(mapStateToProps, {
    newGame,
    insertToken
})(Main)