import ReactDOM from 'react-dom'
import React from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import InputGroup from 'react-bootstrap/InputGroup';
import FormControl from 'react-bootstrap/FormControl';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup'

import { Component } from 'react';

import './App.css';

class TodoList extends Component {

    newTodo = {};

    constructor(props) {
        super(props);

        this.state = {
            todos: []
        }

        this.handleChange = this.handleChange.bind(this);
        this.create = this.create.bind(this);
        this.toggle = this.toggle.bind(this);
    }

    render() {
        return (
            <Container>

                <Row className="justify-content-md-center">
                    <Col md={6} className="text-center">
                        <h2>Todo</h2>
                    </Col>
                </Row>

                <Row className="justify-content-md-center">
                    <Col md={6}>
                        <ButtonGroup className="w-100" vertical>
                            {this.state.todos?.map((todo) =>
                                <ButtonGroup size="lg" key={todo.id}>
                                    <Button onClick={() => this.toggle(todo.id)} className={"todo-item text-left " + (todo.isactive ? "" : "passive")} variant="outline-secondary" key={todo.id}>{todo.title}</Button>
                                    <Button onClick={() => this.delete(todo.id)} variant="outline-secondary" className="delete-button">
                                        Sil
                                    </Button>
                                </ButtonGroup>
                            )}
                        </ButtonGroup>
                    </Col>
                </Row>

                <Row className="justify-content-md-center">
                    <Col md={6}>
                        <InputGroup className="input-group-lg">
                            <FormControl id="newTodo" onChange={this.handleChange} />
                            <InputGroup.Append>
                                <Button variant="outline-secondary" onClick={() => this.create()}>Ekle</Button>
                            </InputGroup.Append>
                        </InputGroup>
                    </Col>
                </Row>

            </Container>
        )
    }

    componentDidMount() {
        this.getTodos();
    }

    handleChange(event) {
        this.newTodo = event.target.value;
    }

    create() {

        fetch('http://localhost:5001/api/v1/todo/create', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: this.newTodo
            }),
        }).then(res => res.json())
            .then((data) => {
                if (data) {
                    this.getTodos();
                }

                document.getElementById('newTodo').value = "";
            })
            .catch(console.log);

    }

    toggle(id) {

        fetch('http://localhost:5001/api/v1/todo/toggle/' + id, {
            method: 'POST',
            headers: {
                'Accept': 'application/json'
            }
        }).then(res => res.json())
            .then((data) => {
                if (data) {
                    this.getTodos();
                }
            })
            .catch(console.log);

    }

    getTodos() {
        fetch('http://localhost:5001/api/v1/todo/list')
            .then(res => res.json())
            .then((data) => {
                this.setState({ todos: data })
            })
            .catch(console.log)
    }

    delete(id) {

        fetch('http://localhost:5001/api/v1/todo/delete/' + id, {
            method: 'DELETE'
        }).then(res => {
            this.getTodos();
        }).catch(console.log);

    }

}

ReactDOM.render(
    <TodoList />, document.getElementById('root')
);

export default TodoList;