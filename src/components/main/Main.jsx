import React, { Component } from 'react';

import './Main.scss';

export default class Main extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <main>
                <h2>FELIX ORTMANN</h2>
                <div className='about__subtitle'>
                    &mdash; Student, Entwickler, Verrückter &mdash;
                </div>
                <div className='about__img'>
                    <img src='img/little_me.png' alt='Fixel' />
                </div>
                <div className='about__description'>
                    <p>
                        Baue gerne verteilte Systeme und mag Herausforderungen. Am liebsten bin ich auf Linux Servern unterwegs,
                        eine Shell ist alles was ich brauche. Handwerke immer wieder in allen möglichen Gebieten wie Web oder Bigdata 
                        und lerne begeistert neue Dinge.
                    </p>
                    <p>
                        Momentan studiere ich an der Uni Hamburg und bin im zweiten Master Semester. Entwickle Software seit 2011, 
                        es begann mit einem Forschungsprojekt zur verteilten Petrinetz Simulation an der Uni. In der Wirtschaft habe ich 
                        mich früh in die Richtung Web bewegt und gleite nun zunehmend in DevOps Bereiche über.
                    </p>
                    <p>
                        Seit Firmengründung bin ich bei <a href='https:joblift.de'>joblift</a> dabei und mache Fullstack. Nebenher bastel ich
                        an einem privaten Opensource Projekt (<a href='https://github.com/0ortmann/wg-tools'>wg-tools</a>), 
                        das bald online geht! Ein paar nette Kleinigkeiten und nützliche Dinge sind auf meinem Github-Account zu finden.
                    </p>
                </div>
            </main>
        );
    }
}